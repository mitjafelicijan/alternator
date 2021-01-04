package main

import "os"

// GetWorkingDirectory ... gets current working directory
func GetWorkingDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}
