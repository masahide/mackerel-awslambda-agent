package store

// Store interface
type Store interface {
	ScanAll(out interface{}) error
	Get(key string, out interface{}) error
	Put(in interface{}) error
}

// ItemNotFoundError error of item not found
//type ItemNotFoundError struct {
//	msg string // description of error
//}
//func (e *ItemNotFoundError) Error() string { return e.msg }
