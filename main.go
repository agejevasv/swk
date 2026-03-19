package main

import (
	"os"

	"github.com/agejevasv/swk/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
