package store

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
)

// LevelStore is the goleveldb storage
type LevelStore struct {
	path string
	db   *leveldb.DB
}

// NewLevelStore returns a new LevelStore
func NewLevelStore(path string) (*LevelStore, error) {
	option := &opt.Options{Compression: opt.SnappyCompression}
	db, err := leveldb.OpenFile(path, option)
	if err != nil {
		return nil, err
	}
	ls := new(LevelStore)
	ls.path = path
	ls.db = db

	return ls, nil
}

// Set implements the Set interface
func (l *LevelStore) Set(key string, data []byte) error {
	return l.db.Put([]byte(key), data, nil)
}

// Get implements the Get interface
func (l *LevelStore) Get(key string) ([]byte, error) {
	data, err := l.db.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	} else {
		return data, err
	}
}

// Del implements the Del interface
func (l *LevelStore) Del(key string) error {
	return l.db.Delete([]byte(key), nil)
}

// Close implements the Close interface
func (l *LevelStore) Close() error {
	err := l.db.Close()
	if err != nil {
		log.Printf("leveldb close error: %s", err)
		return err
	}
	return nil
}

// List implements the List interface
func (l *LevelStore) List(prefix string) ([][]byte, error) {
	var slice *util.Range
	if prefix != "" {
		slice = util.BytesPrefix([]byte(prefix))
	} else {
		slice = nil
	}
	result := make([][]byte, 0)
	iter := l.db.NewIterator(slice, nil)
	for iter.Next() {
		// must copy value from iter.Value()
		value := make([]byte, len(iter.Value()))
		copy(value, iter.Value())
		result = append(result, value)
	}
	iter.Release()
	return result, iter.Error()
}
