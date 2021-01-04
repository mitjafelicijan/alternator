package main

import (
	"fmt"
	"log"
	"net/http"
)

// StartHTTPServer ...
func StartHTTPServer(port int) {

	fs := http.FileServer(http.Dir(fmt.Sprintf("%s/public", GetWorkingDirectory())))
	http.Handle("/", fs)

	log.Println(fmt.Sprintf("Listening on http://0.0.0.0:%d", port))

	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
