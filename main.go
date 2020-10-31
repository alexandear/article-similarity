package main

import (
	"log"

	"github.com/devchallenge/article-similarity/cmd"
)

func main() {
	if err := cmd.ExecuteServer(); err != nil {
		log.Fatal(err)
	}
}
