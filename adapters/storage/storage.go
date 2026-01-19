package storage

// Storage is the interface for the storage layer.
type Storage interface {
	Close() error
}
