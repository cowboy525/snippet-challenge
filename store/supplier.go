package store

import (
	"github.com/ernie-mlg/ErniePJT-main-api-go/store"
)

// SupplierStores structure
type SupplierStores struct {
	space store.SpaceStore
}

// Supplier structure
type Supplier struct {
	stores SupplierStores
}

// NewSupplier : Create new Supplier
func NewSupplier() *Supplier {
	supplier := &Supplier{}

	supplier.CreateStores()

	return supplier
}

// CopyStore : Create new Supplier
func (ss *Supplier) CopyStore() store.Store {
	supplier := &Supplier{}
	supplier.CreateStores()

	return supplier
}

// CreateStores : Create small stores
func (ss *Supplier) CreateStores() {
	ss.stores.space = newSpaceStore(ss)
}

// Space : return SpaceStore instance
func (ss *Supplier) Space() store.SpaceStore {
	return ss.stores.space
}
