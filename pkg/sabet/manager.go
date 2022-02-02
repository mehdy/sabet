package sabet

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/mehdy/sabet/pkg/jobs/meta"
	"sigs.k8s.io/yaml"
)

type Manager struct {
	registry *Registry
	jobs     map[string]meta.Job
	eventsCh chan *meta.Event
	wg       sync.WaitGroup
}

func NewManager() *Manager {
	m := &Manager{
		registry: NewRegistry(),
		jobs:     make(map[string]meta.Job),
		eventsCh: make(chan *meta.Event, 64),
	}
	m.loadConfigs()

	return m
}

func (m *Manager) dispatchEvent(e *meta.Event) {
	m.wg.Add(1)
	m.eventsCh <- e
}

func (m *Manager) handleEvent(e *meta.Event) {
	defer m.wg.Done()

	result, err := e.Job.Execute(e.Payload)
	if err != nil {
		panic(err)
	}

	for _, job := range m.jobs {
		if job.GetRun().SelectorMatch(e.Job.GetLabels()) {
			m.dispatchEvent(&meta.Event{
				Job:     job,
				Payload: result,
			})
		}
	}
}

func (m *Manager) handleEvents() {
	for e := range m.eventsCh {
		go m.handleEvent(e)
	}
}

func (m *Manager) Run() {
	go m.handleEvents()

	for _, job := range m.jobs {
		if job.GetRun().When == meta.RunWhenAlways {
			m.dispatchEvent(&meta.Event{
				Job:     job,
				Payload: nil,
			})
		}
	}

	m.wg.Wait()
}

func (m *Manager) loadConfigs() {
	if err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".yaml") {
			m.loadConfig(path)
		}

		return nil

	}); err != nil {
		panic(err)
	}
}

func (m *Manager) loadConfig(path string) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var dataBuf bytes.Buffer
	tee := io.TeeReader(bytes.NewReader(input), &dataBuf)

	metaBytes, err := ioutil.ReadAll(tee)
	if err != nil {
		panic(err)
	}

	t := &meta.TypeMeta{}
	if err := yaml.Unmarshal(metaBytes, t); err != nil {
		panic(err)
	}

	jobType := m.registry.GetType(t.Type)

	job := reflect.New(jobType.Elem()).Interface().(meta.Job)

	if err := yaml.Unmarshal(dataBuf.Bytes(), job); err != nil {
		panic(err)
	}

	m.jobs[fmt.Sprintf("%s.%s", job.GetType(), job.GetName())] = job
}
