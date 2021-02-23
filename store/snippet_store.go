package store

import "github.com/topoface/snippet-challenge/model"

// SnippetStore structure
type SnippetStore struct {
	*Store
}

func newSnippetStore(Store *Store) *SnippetStore {
	s := &SnippetStore{
		Store,
	}

	return s
}

// Create creates a new snippet
func (ss SnippetStore) Create(snippet *model.Snippet) (*model.Snippet, *model.AppError) {
	return nil, nil
}
