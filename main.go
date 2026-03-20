package main

import (
	"errors"
	"os"

	"github.com/agejevasv/swk/cmd"
	"github.com/agejevasv/swk/internal/ioutil"
)

func main() {
	if err := cmd.Execute(); err != nil {
		var ec ioutil.ExitCoder
		if errors.As(err, &ec) {
			os.Exit(ec.ExitCode())
		}
		os.Exit(2)
	}
}
