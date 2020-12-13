package main

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

var err error
var config *ini.File

func main() {

	init := flag.Bool("init", false, "initializes new page in current directory")
	http := flag.Bool("http", false, "starts HTTP server")
	watch := flag.Bool("watch", false, "watch for file changes and rebuilds HTML")
	build := flag.Bool("build", false, "rebuilds HTML")

	flag.Parse()

	if *init {
		InitializeStaticTemplate()
		os.Exit(0)
	}

	if _, err := os.Stat("config.ini"); os.IsNotExist(err) {
		fmt.Println("Config file does not exist")
		fmt.Println("Use 'staticgen --init' to initialize project")
		os.Exit(1)
	}

	config = ReadConfig()

	defaultTitle := config.Section("content").Key("title").String()
	defaultDescription := config.Section("content").Key("description").String()
	publicFolder := config.Section("server").Key("public").String()
	serverPort := config.Section("server").Key("port").MustInt()

	if *build {
		InitializeMarkdownParser()
		GenerateHTMLFiles(defaultTitle, defaultDescription, publicFolder)
	}

	if *http && !*watch {
		StartHTTPServer(serverPort, publicFolder)
	} else if *http {
		go StartHTTPServer(serverPort, publicFolder)
	}

	if *watch {
		InitializeMarkdownParser()
		GenerateHTMLFiles(defaultTitle, defaultDescription, publicFolder)
		InitializeWatcher(defaultTitle, defaultDescription, publicFolder)
	}

	if !*build && !*http && !*watch && !*init {
		fmt.Println("Try --help option")
	}
}
