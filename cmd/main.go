package cmd

import (
	"log"
	"os"

	"github.com/ernie-mlg/ErniePJT-main-api-go/cmd/commands"
	"github.com/joho/godotenv"
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
