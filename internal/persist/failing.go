package persist

import "errors"

// FailingStore Save 总失败，用于验证 Load 不会激活。
type FailingStore struct {
        Inner Store
        Err   error
}

func (f *FailingStore) Save(id string, raw []byte) error {
        if f.Err != nil {
                return f.Err
        }
        return errors.New("persist: forced failure")
}

func (f *FailingStore) Load(id string) ([]byte, error) {
        if f.Inner != nil {
                return f.Inner.Load(id)
        }
        return nil, ErrNotFound
}

func (f *FailingStore) LoadPath(path string) ([]byte, error) {
        if f.Inner != nil {
                return f.Inner.LoadPath(path)
        }
        return nil, ErrNotFound
}

func (f *FailingStore) List() ([]string, error) {
        if f.Inner != nil {
                return f.Inner.List()
        }
        return nil, nil
}
