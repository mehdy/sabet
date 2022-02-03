package meta

import (
	"bytes"
	"io"
)

type Job interface {
	GetType() string
	GetName() string
	GetLabels() map[string]string
	GetRun() RunSpec

	GetStoreType() string
	SetStore(store Store)

	Execute(input io.Reader) (*bytes.Buffer, error)
}

type Store interface {
	Init() error
	Get(key string) ([]byte, error)
	Put(key string, val []byte) error
	Delete(key string) error
}
