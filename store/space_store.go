package store

// SnippetStore structure
type SnippetStore struct {
	Store
}

func newSnippetStore(Store Store) *SnippetStore {
	s := &SnippetStore{
		Store,
	}

	return s
}

// // Create creates a new snippet
// func (ss SnippetStore) Create(snippet *model.Snippet) (*model.Snippet, *model.AppError) {

// }
