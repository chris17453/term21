package main

import (
	"flag"
	"fmt"
	"term21/term"
)

func usage() {
	fmt.Println("Usage: term21  --theme <theme> <cast.file>")
}

func main() {
	var theme string
	flag.StringVar(&theme, "theme", "Game", "# of iterations")
	flag.Parse()

	var args []string
	args = flag.Args()
	if len(args) <= 0 {
		usage()
		return
	}

	// setup config for terminal
	var t21 term.Term
	var config term.Config
	config.Cast_File = args[0]
	config.Theme = theme

	// init terminal data structure
	if err := t21.Init(config); err != nil {
		fmt.Println(err)
		return
	}
	t21.GifStream()

}
