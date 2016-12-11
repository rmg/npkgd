package main

import (
	"log"
	"net/http"
	"os"
)

// NPMJS is the URL of the official npm registry operated by npm, Inc.
const NPMJS = "https://registry.npmjs.org/"

func main() {
	upstream := NPMJS
	if len(os.Args) > 1 {
		upstream = os.Args[1]
		log.Println("overridding upstream: ", upstream)
	}
	handler := NewNpkgdServer(upstream)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
