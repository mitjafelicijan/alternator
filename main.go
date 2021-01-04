package main

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

var err error
var configFile *ini.File

func main() {

	// path, err := os.Getwd()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(path)

	init := flag.Bool("init", false, "initializes new page in current directory")
	http := flag.Bool("http", false, "starts HTTP server")
	watch := flag.Bool("watch", false, "watch for file changes and rebuilds HTML")
	build := flag.Bool("build", false, "rebuilds HTML")

	flag.Parse()

	if *init {
		InitializeEmptyProject()
		os.Exit(0)
	}

	if _, err := os.Stat("config.ini"); os.IsNotExist(err) {
		fmt.Println("Config file does not exist")
		fmt.Println("Use 'alternator --init' to initialize project")
		fmt.Println("More info at https://github.com/mitjafelicijan/alternator")
		os.Exit(1)
	}

	configFile = ReadConfig()

	defaultTitle := configFile.Section("content").Key("title").String()
	defaultDescription := configFile.Section("content").Key("description").String()
	publicFolder := configFile.Section("generator").Key("public").String()
	serverPort := configFile.Section("server").Key("port").MustInt()

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
