package sabet

import (
	"reflect"

	"github.com/mehdy/sabet/pkg/jobs/meta"
	"github.com/mehdy/sabet/pkg/jobs/rss"
	"github.com/mehdy/sabet/pkg/jobs/telegram"
	"github.com/mehdy/sabet/pkg/stores/fs"
)

type Registry struct {
	jobTypes   map[string]reflect.Type
	storeTypes map[string]reflect.Type
}

func NewRegistry() *Registry {
	r := &Registry{
		jobTypes:   make(map[string]reflect.Type),
		storeTypes: make(map[string]reflect.Type),
	}
	r.registerAll()

	return r
}

func (r *Registry) registerAll() {
	r.RegisterJobType("RSS", &rss.Job{})
	r.RegisterJobType("Telegram", &telegram.Job{})

	r.RegisterStoreType("fs", &fs.FS{})
}

func (r *Registry) RegisterJobType(name string, t meta.Job) {
	r.jobTypes[name] = reflect.TypeOf(t)
}

func (r *Registry) RegisterStoreType(name string, t meta.Store) {
	r.storeTypes[name] = reflect.TypeOf(t)
}

func (r *Registry) GetStoreType(name string) reflect.Type {
	return r.storeTypes[name]
}

func (r *Registry) GetJobType(name string) reflect.Type {
	return r.jobTypes[name]
}
