package store

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Slug  string `gorm:"uniqueIndex:idx_slug"`
	Title string
	Body  string
}

var db *gorm.DB

func Open() *gorm.DB {
	var err error
	db, err = gorm.Open(sqlite.Open("data/sqlite.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	return db
}

func Migrate() *gorm.DB {
	if db == nil {
		Open()
	}
	db.AutoMigrate(&Post{})
	return db
}

/*
================
Helper Functions
================
*/
func CreatePost(title string, body string) (Post, error) {
	token, err := generateHash(32)
	if err != nil {
		return Post{}, err
	}

	newPost := Post{Title: title, Body: body, Slug: token}

	if res := db.Create(&newPost); res.Error != nil {
		return Post{}, res.Error
	}

	return newPost, nil
}

// ErrNotFound is returned by lookups that match no row.
var ErrNotFound = errors.New("not found")

// ListPosts returns all posts, newest first.
func ListPosts() ([]Post, error) {
	var posts []Post
	if res := db.Order("created_at desc, id desc").Find(&posts); res.Error != nil {
		return nil, res.Error
	}
	return posts, nil
}

// GetPostBySlug returns the post with the given slug, or ErrNotFound.
func GetPostBySlug(slug string) (Post, error) {
	var posts []Post
	if res := db.Where("slug = ?", slug).Limit(1).Find(&posts); res.Error != nil {
		return Post{}, res.Error
	}
	if len(posts) == 0 {
		return Post{}, ErrNotFound
	}
	return posts[0], nil
}

/*
=================
Hashing function
=================
*/
func generateHash(s int) (string, error) {
	// b, err := generateRandBytes(s)

	// return base64.URLEncoding.EncodeToString(b), err
	b := make([]byte, s)
	_, err := rand.Read(b)

	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(b), err
}
