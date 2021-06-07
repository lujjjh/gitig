package main

import (
	"log"
	"net/http"

	"github.com/lujjjh/gitig/git"
)

func main() {
	if err := http.ListenAndServe(":3000", new(git.SmartHTTPHandler)); err != nil {
		log.Fatal(err)
	}
}
