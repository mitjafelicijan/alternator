package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/packr/v2/file"
)

var content string

var configFile string = `[server]
# port used by HTTP server
port = 8080

# location of generated html files
public = ./public

[content]
# title displayed on index page, otherwise post title is used
title = My website

# descriptions displayed on index page, otherwise post description is used
description = Default description`

// InitializeStaticTemplate ...
func InitializeStaticTemplate() {
	root := "./"
	// root := "./demo/"

	fmt.Println("Initializing new project")

	openVirtualDirectoryAndCopy("./template", fmt.Sprintf("%s%s", root, "template"))
	openVirtualDirectoryAndCopy("./posts", fmt.Sprintf("%s%s", root, "posts"))
	openVirtualDirectoryAndCopy("./assets", fmt.Sprintf("%s%s", root, "assets"))

	createDirectoryIfNotExists(fmt.Sprintf("%s%s", root, "public"))

	writeContentToFile(configFile, fmt.Sprintf("%s%s", root, "config.ini"))
}

func openVirtualDirectoryAndCopy(virtualDirectory string, realDirectory string) {
	createDirectoryIfNotExists(realDirectory)

	box := packr.NewBox(virtualDirectory)
	err := box.Walk(func(path string, f file.File) error {
		content, _ = box.FindString(path)
		writeContentToFile(content, fmt.Sprintf("%s/%s", realDirectory, path))
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func writeContentToFile(content string, path string) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		log.Println(err)
	}
}

func createDirectoryIfNotExists(directory string) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, 0700)
	}
}
