package main

import (
	"os"

	"github.com/antonpodkur/remember/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
