package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

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

// InitializeEmptyProject ...
func InitializeEmptyProject() {
	fmt.Println("Initializing new project")

	filepaths := AssetNames()

	for _, filepath := range filepaths {
		dir := path.Dir(filepath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}

		data, err := Asset(filepath)
		if err != nil {
			panic(err)
		}

		writeContentToFile(string(data), filepath)
	}

	createDirectoryIfNotExists(fmt.Sprintf("%s%s", "./", "public"))
}
