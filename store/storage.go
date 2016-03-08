package store

const (
	errNotExisted     string = "Data Not Existed"
	errModeNotMatched string = "Storage Mode Not Matched"
)

// Storage is the storage interface based on key-value pattern
type Storage interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	Del(key string) error
	Close() error
	List(prefix string) ([][]byte, error)
}
