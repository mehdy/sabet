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
	Execute(input io.Reader) (*bytes.Buffer, error)
}
