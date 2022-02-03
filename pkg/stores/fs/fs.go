package fs

import (
	"io/ioutil"
	"os"
	"path"
)

type FS struct {
	Path string `json:"path,omitempty"`
}

func (f *FS) Init() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	f.Path = path.Join(wd, f.Path)

	return os.MkdirAll(f.Path, 0755)
}

func (f *FS) Get(key string) ([]byte, error) {
	return ioutil.ReadFile(path.Join(f.Path, key))
}

func (f *FS) Put(key string, val []byte) error {
	return os.WriteFile(path.Join(f.Path, key), val, 0644)
}

func (f *FS) Delete(key string) error {
	return os.Remove(path.Join(f.Path, key))
}
