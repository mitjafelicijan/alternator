package main

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

var config Config

// Config ...
type Config struct {
	Server struct {
		Port int
	}
	Generator struct {
		Public string
	}
	Content struct {
		Title       string
		Description string
	}
	RSS struct {
		Domain string
		Author string
		Email  string
	}
}

// ReadConfig ...
func ReadConfig() *ini.File {
	configFile, err = ini.Load(fmt.Sprintf("%s/config.ini", GetWorkingDirectory()))
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	return configFile
}
