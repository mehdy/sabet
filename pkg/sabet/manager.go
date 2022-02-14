package sabet

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"sync"

	"github.com/mehdy/sabet/pkg/jobs/meta"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
)

type Manager struct {
	registry *Registry
	jobs     map[string]meta.Job
	eventsCh chan *meta.Event
	wg       sync.WaitGroup
}

func NewManager(configDirs []string) *Manager {
	m := &Manager{
		registry: NewRegistry(),
		jobs:     make(map[string]meta.Job),
		eventsCh: make(chan *meta.Event, 64),
	}
	m.loadConfigs(configDirs)

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
		logrus.WithError(err).WithField("job", e.Job.GetName()).Error("Job failed")
		return
	}

	logrus.WithField("job", e.Job.GetName()).Info("Job finished")

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

func (m *Manager) loadConfigs(configDirs []string) {
	for _, dir := range configDirs {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			logrus.WithError(err).Errorf("Skipped reading configs inside %q", dir)
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
				continue
			}

			if err := m.loadConfig(file.Name()); err != nil {
				logrus.WithError(err).WithField("file", file.Name()).Error("Error while loading config")
			}
		}
	}
}

func (m *Manager) loadConfig(path string) error {
	logrus.WithField("file", path).Info("Loading config")

	input, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	job, err := m.loadJobMeta(input)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(input, job); err != nil {
		return err
	}

	if err := m.loadStore(job, input); err != nil {
		return err
	}

	if err := job.Init(); err != nil {
		return err
	}

	m.jobs[fmt.Sprintf("%s.%s", job.GetType(), job.GetName())] = job

	return nil
}

func (m *Manager) loadJobMeta(buf []byte) (meta.Job, error) {
	t := &meta.TypeMeta{}
	if err := yaml.Unmarshal(buf, t); err != nil {
		return nil, err
	}

	jobType := m.registry.GetJobType(t.Type)
	if jobType == nil {
		return nil, fmt.Errorf("Unknown job type %q", t.Type)
	}

	return reflect.New(jobType.Elem()).Interface().(meta.Job), nil
}

type StoreLoader struct {
	meta.Store `json:"store,omitempty"`
}

func (m *Manager) loadStore(job meta.Job, buf []byte) error {
	if job.GetStoreType() == "" {
		return nil
	}
	storeType := m.registry.GetStoreType(job.GetStoreType())

	store := reflect.New(storeType.Elem()).Interface().(meta.Store)

	sl := &StoreLoader{store}
	if err := yaml.Unmarshal(buf, sl); err != nil {
		return err
	}

	if err := store.Init(); err != nil {
		return err
	}

	job.SetStore(store)

	return nil
}
