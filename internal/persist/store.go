package persist

import "errors"

var ErrNotFound = errors.New("persist: not found")

type Store interface {
        Save(id string, raw []byte) error
        Load(id string) ([]byte, error)
        LoadPath(path string) ([]byte, error)
        List() ([]string, error)
}
