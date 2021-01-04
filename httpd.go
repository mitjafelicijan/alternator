package main

import (
	"fmt"
	"log"
	"net/http"
)

// StartHTTPServer ...
func StartHTTPServer(port int, dir string) {
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	log.Println(fmt.Sprintf("Listening on http://0.0.0.0:%d", port))

	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
