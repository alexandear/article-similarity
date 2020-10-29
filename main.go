package main

import (
	"log"

	cmder "github.com/yaegashi/cobra-cmder"

	"github.com/devchallenge/article-similarity/cmd"
)

func main() {
	app := &cmd.App{}
	if err := cmder.Cmd(app).Execute(); err != nil {
		log.Fatal(err)
	}
}
