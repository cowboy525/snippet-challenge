package testutils

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/topoface/snippet-challenge/utils/fileutils"
)

func ReadTestFile(name string) ([]byte, error) {
	path, _ := fileutils.FindDir("tests")
	file, err := os.Open(filepath.Join(path, name))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := &bytes.Buffer{}
	if _, err := io.Copy(data, file); err != nil {
		return nil, err
	} else {
		return data.Bytes(), nil
	}
}

// LoadEnv loads env vars from .env
func LoadEnv() {
	srcPath := fileutils.FindFile(".env")

	err := godotenv.Load(srcPath)
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
		os.Exit(-1)
	}
}
