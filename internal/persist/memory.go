package persist

import (
        "os"
        "sync"
)

type MemoryStore struct {
        mu   sync.RWMutex
        data map[string][]byte
}

func NewMemoryStore() *MemoryStore {
        return &MemoryStore{data: make(map[string][]byte)}
}

func (m *MemoryStore) Save(id string, raw []byte) error {
        m.mu.Lock()
        defer m.mu.Unlock()
        cp := make([]byte, len(raw))
        copy(cp, raw)
        m.data[id] = cp
        return nil
}

func (m *MemoryStore) Load(id string) ([]byte, error) {
        m.mu.RLock()
        defer m.mu.RUnlock()
        b, ok := m.data[id]
        if !ok {
                return nil, ErrNotFound
        }
        cp := make([]byte, len(b))
        copy(cp, b)
        return cp, nil
}

func (m *MemoryStore) LoadPath(path string) ([]byte, error) { return os.ReadFile(path) }

func (m *MemoryStore) List() ([]string, error) {
        m.mu.RLock()
        defer m.mu.RUnlock()
        out := make([]string, 0, len(m.data))
        for k := range m.data {
                out = append(out, k)
        }
        return out, nil
}
