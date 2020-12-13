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

	"github.com/otiai10/copy"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var md goldmark.Markdown

// Post ...
type Post struct {
	File        string
	Title       string
	Description string
	Slug        string
	Content     string
	Tags        interface{}
}

// Index ...
type Index struct {
	Title       string
	Description string
	Posts       []Post
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

// GenerateHTMLFiles ...
func GenerateHTMLFiles(defaultTitle string, defaultDescription string, publicFolder string) {
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

		post := Post{
			File:        file,
			Title:       meta["Title"].(string),
			Description: meta["Description"].(string),
			Slug:        fmt.Sprintf("%s.html", meta["Slug"].(string)),
			Content:     html,
			Tags:        meta["Tags"],
		}

		// tags = append(tags, "item")

		fmt.Println(reflect.ValueOf(meta["Tags"]).Kind())

		for key, value := range meta["Tags"] {
			fmt.Println(key, value)
		}

		//fmt.Println(meta["Tags"].([]string))
		// fmt.Println(meta["Tags"])

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

	fmt.Println(tags)

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
