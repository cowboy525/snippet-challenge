package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/topoface/snippet-challenge/model"
	"github.com/topoface/snippet-challenge/utils/fileutils"
)

// Date Event part
func ReadSnippets() []*model.Snippet {
	pth := path.Join(fileutils.RootDir(), "snippets.txt")

	bytes, err := ioutil.ReadFile(pth)
	if err != nil {
		return []*model.Snippet{}
	}

	var data []*model.Snippet
	json.Unmarshal(bytes, &data)
	return data
}

func WriteSnippets(data []*model.Snippet) {
	pth := path.Join(fileutils.RootDir(), "snippets.txt")

	// If the file doesn't exist, create it, or append to the file
	out, err := os.OpenFile(pth, os.O_APPEND|os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	b, _ := json.Marshal(data)
	out.Write(b)
}

func (a *App) CreateSnippet(request *model.SnippetRequest) (*model.Snippet, *model.AppError) {
	url := fmt.Sprintf("%v/snippets/%v", os.Getenv("HOST_URL"), request.Name)
	newSnippet := &model.Snippet{
		Name:      request.Name,
		ExpiresAt: time.Now().Add(time.Second * time.Duration(request.ExpiresIn)),
		Body:      request.Body,
		URL:       url,
	}

	data := ReadSnippets()

	sameNameExists := false
	for _, s := range data {
		if s.Name == newSnippet.Name {
			s.ExpiresAt = newSnippet.ExpiresAt
			s.Body = newSnippet.Body
			sameNameExists = true
		}
	}

	if !sameNameExists {
		data = append(data, newSnippet)
	}

	WriteSnippets(data)

	return newSnippet, nil
}

func (a *App) GetSnippet(name string) (*model.Snippet, *model.AppError) {
	data := ReadSnippets()

	for _, s := range data {
		if s.Name == name {
			if s.ExpiresAt.Before(time.Now()) {
				break
			}
			s.ExpiresAt = s.ExpiresAt.Add(time.Second * 30)
			WriteSnippets(data)
			return s, nil
		}
	}

	return nil, model.NotFoundError("", "")
}
