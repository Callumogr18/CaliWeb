package server

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"mysite/src/store"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/style.css
var styleCSS []byte

// pages holds one template set per page (layout + that page), parsed once at
// startup. A template error panics here, so the server fails at boot and never
// at request time.
var pages = mustParsePages("home", "thoughts", "post", "projects", "notfound", "error")

func mustParsePages(names ...string) map[string]*template.Template {
	m := make(map[string]*template.Template, len(names))
	for _, name := range names {
		m[name] = template.Must(template.ParseFS(templateFS, "templates/layout.html", "templates/"+name+".html"))
	}
	return m
}

// Template data. Every page embeds Page, which layout.html reads.
// See docs/ui-spec.md, "Template data contracts".
type Page struct {
	Title  string
	Active string // "home" | "thoughts" | "projects" | "" (nav highlight)
	Year   int
}

type HomePage struct {
	Page
	Latest *store.Post // nil when there are no posts
}

type ThoughtsPage struct {
	Page
	Posts []store.Post // nil or empty is fine
}

type PostPage struct {
	Page
	Post       store.Post
	Paragraphs []string
}

func newPage(title, active string) Page {
	return Page{Title: title, Active: active, Year: time.Now().Year()}
}

// render executes a page into a buffer first and only then writes the status
// and body, so a template error yields a clean 500 instead of a half page.
func render(w http.ResponseWriter, status int, name string, data any) {
	t, ok := pages[name]
	if !ok {
		log.Printf("render: unknown page %q", name)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "layout", data); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; base-uri 'none'; form-action 'none'")
	w.WriteHeader(status)
	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Printf("write %s: %v", name, err)
	}
}

var blankLine = regexp.MustCompile(`\n\s*\n`)

// paragraphs splits a post body on blank lines. The result is plain text that
// html/template escapesthe body is never treated as HTML.
func paragraphs(body string) []string {
	body = strings.ReplaceAll(body, "\r\n", "\n")
	var out []string
	for _, p := range blankLine.Split(body, -1) {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// styleETag identifies the embedded stylesheet, so it changes whenever the CSS does.
var styleETag = func() string {
	sum := sha256.Sum256(styleCSS)
	return `"` + hex.EncodeToString(sum[:8]) + `"`
}()

// styleSheet serves the embedded stylesheet at GET /static/style.css.
//
// The CSS is baked into the binary and changes on every deploy, so it must not be
// cached blindly: "no-cache" makes the browser revalidate on each load, and the
// ETag lets the server answer 304 Not Modified when nothing has changed.
func styleSheet(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/css; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Cache-Control", "no-cache")
	h.Set("ETag", styleETag)
	http.ServeContent(w, r, "style.css", time.Time{}, bytes.NewReader(styleCSS))
}
