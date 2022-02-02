package sabet

import (
	"reflect"

	"github.com/mehdy/sabet/pkg/jobs/meta"
)

type Registry struct {
	nameToType map[string]reflect.Type
}

func NewRegistry() *Registry {
	r := &Registry{
		nameToType: make(map[string]reflect.Type),
	}
	r.registerAll()

	return r
}

func (r *Registry) Register(name string, t meta.Job) {
	r.nameToType[name] = reflect.TypeOf(t)
}

func (r *Registry) GetType(name string) reflect.Type {
	return r.nameToType[name]
}
