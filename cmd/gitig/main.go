package main

import (
	"log"
	"net/http"

	"github.com/lujjjh/gitig/git/smarthttp"
)

func main() {
	if err := http.ListenAndServe(":3000", new(smarthttp.Handler)); err != nil {
		log.Fatal(err)
	}
}
