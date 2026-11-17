# Personal site: Go + net/http

## Context
Empty project, Go 1.26 installed. Goal: a personal site that starts as a simple server-rendered site and can grow into a blog, then possibly a community space. Keep dependencies near zero; Go 1.22+ `ServeMux` supports method and path-wildcard routing, so no router library is needed.

## The stack
| Layer | Choice | Why |
|---|---|---|
| Server | `net/http` + stdlib `ServeMux` (`GET /posts/{slug}`) | No framework needed |
| Content (phase 1) | Markdown files with front matter in `content/posts/`, parsed at startup with `goldmark` + `yaml` | Write posts in an editor, git is the CMS |
| Responses | JSON via `encoding/json` | API-only backend; no frontend in this repo |
| Storage (phase 2+) | SQLite (`modernc.org/sqlite`, pure Go) + `sqlc` or plain `database/sql` | Zero ops, one file, enough for a community |
| Auth (phase 3) | Cookie sessions (`gorilla/sessions` or hand-rolled), `bcrypt`/`argon2`; or OAuth (GitHub) | Only needed for community features |
| Middleware | Hand-written: logging, recover, gzip, security headers | Small `func(http.Handler) http.Handler` chain |
| Config / logging | Env vars, `log/slog` | Stdlib |
| Deploy | Single static binary or Docker image; Fly.io / small VPS; Caddy in front for TLS | Cheap and simple |

## Layout
```
website/
  src/
    cmd/main.go            # entry point: builds router, listens on :8080
    server/
      server.go            # NewRouter: every route
      handlers.go          # JSON handlers; notFound is the catch-all "/" handler
    store/db_helper.go     # gorm/SQLite models and helpers
  docs/api.md              # endpoint list (source of truth for routes)
  ai/knowledge/            # AGENT.md (this file)
  go.mod                   # module "mysite"
```

Conventions:
- The frontend (HTML, CSS, JS, htmx, node packages) was removed. This repo is an API only.
- Handlers respond with JSON; errors are `{"error": "..."}`.
- `docs/api.md` lists the routes and `server.go` registers them.

## Phases
1. **Static-ish site (start here):** `/`, `/about`, `/posts`, `/posts/{slug}`; templates and static files embedded; posts read from markdown at startup; `/healthz`; graceful shutdown with timeouts set on `http.Server`.
2. **Blog polish:** RSS/Atom feed, tags, syntax highlighting (`goldmark-highlighting`), sitemap, OpenGraph tags, draft flag, optional hot reload in dev (`air`).
3. **Dynamic data:** SQLite, migrations (`goose`), admin-only post editor or comments.
4. **Community:** accounts, sessions, CSRF protection, rate limiting, moderation, email (SMTP/Resend) for verification.

## Phase 1 implementation steps
1. `go mod init` (module path to be chosen, e.g. `github.com/<user>/personal`).
2. `cmd/web/main.go`: `http.Server` with Read/Write/Idle timeouts, `ServeMux`, signal-based shutdown.
3. `web/` embed package exposing `embed.FS`; template cache built at startup (base layout + per-page).
4. `internal/posts`: parse front matter (title, date, slug, draft), render markdown to `template.HTML`, sort by date.
5. Handlers + routes; 404 page; logging/recover middleware.
6. Minimal CSS, one sample post.

## Verification
- `go run ./cmd/web` then curl/browser `/`, `/posts`, `/posts/hello-world`, an unknown path (404), and `/healthz`.
- `go vet ./... && go test ./...` (add a test for post parsing and a handler test with `httptest`).
- `go build` produces one binary that runs from another directory (confirms embed works).
