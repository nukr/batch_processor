package db

// DB ...
type DB interface {
	Update(selector, document []byte) (int32, error)
}
