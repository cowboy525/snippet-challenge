package store

// Store structure
type Store struct {
	snippet *SnippetStore
}

// NewStore : Create new Supplier
func NewStore() *Store {
	store := &Store{}

	store.CreateStores()

	return store
}

func (ss *Store) CreateStores() {
	// ss.stores.snippet = newSnippetStore(ss)
}
