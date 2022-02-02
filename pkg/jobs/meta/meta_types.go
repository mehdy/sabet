package meta

import "io"

const (
	RunWhenAlways = "always"
)

// TypeMeta describes an individual job in an configuration
// with strings representing the type of the job and its version.
type TypeMeta struct {
	Type string `json:"type,omitempty"`
}

func (t TypeMeta) GetType() string {
	return t.Type
}

// Job is metadata that all jobs must have.
type JobMeta struct {
	Name   string            `json:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
	Run    RunSpec           `json:"run,omitempty"`
}

func (j JobMeta) GetName() string {
	return j.Name
}

func (j JobMeta) GetLabels() map[string]string {
	return j.Labels
}

func (j JobMeta) GetRun() RunSpec {
	return j.Run
}

// RunSpec describes when and if a job should be executed.
type RunSpec struct {
	When     string            `json:"when,omitempty"`
	If       string            `json:"if,omitempty"`
	Selector map[string]string `json:"selector,omitempty"`
}

func (r RunSpec) SelectorMatch(m map[string]string) bool {
	for k, v := range m {
		if sv, ok := r.Selector[k]; !ok || sv != v {
			return false
		}
	}
	return true
}

type Event struct {
	Job     Job
	Payload io.Reader
}
