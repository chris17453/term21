package main

import (
	"fmt"
	"os"
	"term21/term"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name!")
		return
	}

	var t21 term.Term

	// setup config for terminal
	var config term.Config
	config.Cast_File = os.Args[1]
	config.Font_File = "/home/nd/repos/ttygif-assets/fonts/fd/Bm437_AMI_BIOS.fd"

	// init terminal data structure
	if err := t21.Init(config); err != nil {
		fmt.Println(err)
		return
	}

}
