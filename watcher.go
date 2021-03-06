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
func InitializeWatcher() {
	directory := fmt.Sprintf("%s/posts", GetWorkingDirectory())

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
				GenerateHTMLFiles(configFile)

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
