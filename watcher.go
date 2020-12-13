package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

// InitializeWatcher ...
func InitializeWatcher(defaultTitle string, defaultDescription string, publicFolder string) {
	directory := "./posts"

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk(directory, watchDir); err != nil {
		fmt.Println("Error", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case _ = <-watcher.Events:
				GenerateHTMLFiles(defaultTitle, defaultDescription, publicFolder)
				// log.Printf("Event %#v\n", event)

			case err := <-watcher.Errors:
				log.Println("Error", err)
			}
		}
	}()

	<-done
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}
