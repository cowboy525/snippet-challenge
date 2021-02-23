package cmd

import (
	"os"

	"github.com/topoface/snippet-challenge/cmd/commands"
)

func Execute() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
