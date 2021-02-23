package store

// Store interface
type Store interface {
	Space() SpaceStore
}

// SpaceStore interface
type SpaceStore interface {
	// Create(space *model.Space) (*model.Space, *model.AppError)
	// Update(space *model.Space) (*model.Space, *model.AppError)
	// Get(id uint64) (*model.Space, *model.AppError)
}
