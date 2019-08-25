package store

// Store interface
type Store interface {
	ScanAll(out interface{}) error
	Get(key string, out interface{}) error
	Put(in interface{}) error
}
