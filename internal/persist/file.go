package persist

import (
        "os"
        "path/filepath"
        "sync"
)

type FileStore struct {
        dir string
        mu  sync.Mutex
}

func NewFileStore(dir string) *FileStore {
        _ = os.MkdirAll(dir, 0o755)
        return &FileStore{dir: dir}
}

func (f *FileStore) path(id string) string {
        return filepath.Join(f.dir, id+".json")
}

func (f *FileStore) Save(id string, raw []byte) error {
        f.mu.Lock()
        defer f.mu.Unlock()
        tmp := f.path(id) + ".tmp"
        if err := os.WriteFile(tmp, raw, 0o644); err != nil {
                return err
        }
        return os.Rename(tmp, f.path(id))
}

func (f *FileStore) Load(id string) ([]byte, error) {
        b, err := os.ReadFile(f.path(id))
        if err != nil {
                if os.IsNotExist(err) {
                        return nil, ErrNotFound
                }
                return nil, err
        }
        return b, nil
}

func (f *FileStore) LoadPath(path string) ([]byte, error) { return os.ReadFile(path) }

func (f *FileStore) List() ([]string, error) {
        entries, err := os.ReadDir(f.dir)
        if err != nil {
                return nil, err
        }
        var out []string
        for _, e := range entries {
                if e.IsDir() {
                        continue
                }
                name := e.Name()
                if filepath.Ext(name) == ".json" {
                        out = append(out, name[:len(name)-5])
                }
        }
        return out, nil
}
