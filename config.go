package main

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

// ReadConfig ...
func ReadConfig() *ini.File {
	config, err = ini.Load("./config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	return config
}
