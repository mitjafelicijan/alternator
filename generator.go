package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"text/template"
	"time"

	"github.com/gorilla/feeds"
	"github.com/otiai10/copy"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/ini.v1"
)

var md goldmark.Markdown

// Post ...
type Post struct {
	File        string
	Title       string
	Description string
	Slug        string
	Created     string
	Content     string
	Listing     bool
	Tags        interface{}
}

// Index ...
type Index struct {
	Title       string
	Description string
	Posts       []Post
	Tags        interface{}
}

// Tag ...
type Tag struct {
	Title       string
	Tag         string
	Description string
	Posts       []Post
	Tags        interface{}
}

// InitializeMarkdownParser ...
func InitializeMarkdownParser() {
	md = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithBlockParsers(),
			parser.WithInlineParsers(),
			parser.WithParagraphTransformers(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
}

// CleanPublicDirectory ...
func CleanPublicDirectory(directory string) {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if path != directory {
			log.Println(fmt.Sprintf("Removing %s", path))
			os.Remove(path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

// ListMarkdownFiles ... retrieves all markdown files from specific folder
func ListMarkdownFiles(root string) ([]string, error) {
	var matches []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if matched, err := filepath.Match("*.md", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return matches, nil
}

// ConvertMarkdownFileToHTML ... converts md to html and extracts meta title block
func ConvertMarkdownFileToHTML(filepath string) (string, map[string]interface{}, error) {
	source, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Print(err)
		return "", nil, err
	}

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert(source, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	meta := meta.Get(context)

	return string(buf.String()), meta, nil
}

// FindItemInSlice ...
func FindItemInSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// SliceExists ..
func SliceExists(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("SliceExists() given a non-slice type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

// GenerateHTMLFiles ...
func GenerateHTMLFiles(configFile *ini.File) {
	publicFolder := fmt.Sprintf("%s/public", GetWorkingDirectory())

	CleanPublicDirectory(publicFolder)

	files, err := ListMarkdownFiles("./posts")
	if err != nil {
		log.Fatal(err)
	}

	posts := []Post{}

	var tags []string

	tpl, err := template.ParseGlob("./template/*.html")
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		html, meta, _ := ConvertMarkdownFileToHTML(file)

		if meta["Listing"].(bool) {
			post := Post{
				File:        file,
				Title:       meta["Title"].(string),
				Description: meta["Description"].(string),
				Slug:        fmt.Sprintf("%s.html", meta["Slug"].(string)),
				Content:     html,
				Created:     meta["Created"].(string),
				Listing:     meta["Listing"].(bool),
				Tags:        meta["Tags"],
			}

			// append tag to all tags if not there already
			for _, tag := range meta["Tags"].([]interface{}) {
				_, found := FindItemInSlice(tags, tag.(string))
				if !found {
					fmt.Println("Value not found in slice")
					tags = append(tags, tag.(string))
				}
			}

			posts = append(posts, post)

			log.Println(fmt.Sprintf("Generating %s.html file ... ", post.Slug))
			output, err := os.Create(fmt.Sprintf("%s/%s", publicFolder, post.Slug))
			if err != nil {
				log.Println("Create file: ", err)
				return
			}

			err = tpl.ExecuteTemplate(output, "post.html", post)
			if err != nil {
				panic(err)
			}

			output.Close()
		}

	}

	defaultTitle := configFile.Section("content").Key("title").String()
	defaultDescription := configFile.Section("content").Key("description").String()

	// generate tag index and tag listings

	_ = os.Mkdir(fmt.Sprintf("%s/tags", publicFolder), 0777)

	for _, tag := range tags {
		log.Println(fmt.Sprintf("Generating tags/%s.html file ... ", tag))

		output, err := os.Create(fmt.Sprintf("%s/tags/%s.html", publicFolder, tag))
		if err != nil {
			log.Println("Create file: ", err)
			return
		}

		tagPosts := []Post{}

		for _, post := range posts {
			if SliceExists(post.Tags, tag) {
				tagPosts = append(tagPosts, post)
			}
		}

		tag := Tag{
			Title:       defaultTitle,
			Tag:         tag,
			Description: defaultDescription,
			Posts:       tagPosts,
			Tags:        tags,
		}

		err = tpl.ExecuteTemplate(output, "tag.html", tag)
		if err != nil {
			panic(err)
		}

		output.Close()

	}

	// generate rss feed

	domain := configFile.Section("rss").Key("domain").String()
	author := configFile.Section("rss").Key("author").String()
	email := configFile.Section("rss").Key("email").String()

	now := time.Now()
	feed := &feeds.Feed{
		Title:       defaultTitle,
		Link:        &feeds.Link{Href: domain},
		Description: defaultDescription,
		Author:      &feeds.Author{Name: author, Email: email},
		Created:     now,
	}

	feed.Items = []*feeds.Item{}
	for _, post := range posts {
		feedItem := &feeds.Item{
			Title:       post.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/%s", domain, post.Slug)},
			Description: post.Description,
			Content:     post.Content,
			Author:      &feeds.Author{Name: author, Email: email},
			Created:     now,
		}

		feed.Items = append(feed.Items, feedItem)
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}
	writeContentToFile(string(rss), fmt.Sprintf("%s/feed.rss", publicFolder))

	json, err := feed.ToJSON()
	if err != nil {
		log.Fatal(err)
	}
	writeContentToFile(string(json), fmt.Sprintf("%s/feed.json", publicFolder))

	// generate index.html file

	log.Println("Generating index.html file ... ")
	output, err := os.Create(fmt.Sprintf("%s/index.html", publicFolder))
	if err != nil {
		log.Println("Create file: ", err)
		return
	}

	index := Index{
		Title:       defaultTitle,
		Description: defaultDescription,
		Posts:       posts,
		Tags:        tags,
	}

	err = tpl.ExecuteTemplate(output, "index.html", index)
	if err != nil {
		panic(err)
	}

	output.Close()

	CopyFile("./template/script.js", fmt.Sprintf("%s/script.js", publicFolder))
	CopyFile("./template/style.css", fmt.Sprintf("%s/style.css", publicFolder))
	CopyFile("./template/favicon.ico", fmt.Sprintf("%s/favicon.ico", publicFolder))

	CopyAssets(publicFolder)

	fmt.Println()
}

// CopyAssets ...
func CopyAssets(publicFolder string) {
	log.Println("Copying assets to public folder")

	err = copy.Copy("./assets", fmt.Sprintf("%s/assets", publicFolder))
	if err != nil {
		panic(err)
	}
}

// CopyFile ...
func CopyFile(src, dst string) error {
	log.Println(fmt.Sprintf("Copying %s file ... ", src))

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
