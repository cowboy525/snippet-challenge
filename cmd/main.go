package cmd

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/topoface/snippet-challenge/cmd/commands"
)

func Execute() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(-1)
	}

	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
