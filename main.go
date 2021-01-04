package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"gopkg.in/ini.v1"
)

var err error
var configFile *ini.File
var ldVersion = "development"

func main() {
	init := flag.Bool("init", false, "initializes new page in current directory")
	http := flag.Bool("http", false, "starts HTTP server")
	watch := flag.Bool("watch", false, "watch for file changes and rebuilds HTML")
	build := flag.Bool("build", false, "rebuilds HTML")
	version := flag.Bool("version", false, "rebuilds HTML")

	flag.Parse()

	if *init {
		InitializeEmptyProject()
		os.Exit(0)
	}

	if *version {
		fmt.Println("Version", ldVersion)
		os.Exit(0)
	}

	if _, err := os.Stat("config.ini"); os.IsNotExist(err) {
		color.Red("Error: Config file does not exist!")
		fmt.Println("Use `alternator --init` to initialize project.")
		fmt.Println("More info at https://github.com/mitjafelicijan/alternator.")
		os.Exit(1)
	}

	configFile = ReadConfig()
	serverPort := configFile.Section("server").Key("port").MustInt()

	if *build {
		InitializeMarkdownParser()
		GenerateHTMLFiles(configFile)
	}

	if *http && !*watch {
		StartHTTPServer(serverPort)
	} else if *http {
		go StartHTTPServer(serverPort)
	}

	if *watch {
		InitializeMarkdownParser()
		GenerateHTMLFiles(configFile)
		InitializeWatcher()
	}

	if !*build && !*http && !*watch && !*init {
		fmt.Println("Try --help option")
	}
}
