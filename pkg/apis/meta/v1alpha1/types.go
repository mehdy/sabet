package v1alpha1

import "fmt"

// TypeMeta describes an individual job in an configuration
// with strings representing the type of the job and its API schema version.
// Structures that are versioned or persisted should inline TypeMeta.
type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

func (t TypeMeta) GetKind() string {
	return t.Kind
}

func (t TypeMeta) GetAPIVersion() string {
	return t.APIVersion
}

func (t TypeMeta) GetVersionKind() string {
	return fmt.Sprintf("%s.%s", t.Kind, t.APIVersion)
}

// JobMeta is metadata that all persisted jobs must have.
type JobMeta struct {
	Name   string            `json:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
	Run    RunSpec           `json:"run,omitempty"`
}

func (j *JobMeta) GetName() string {
	return j.Name
}

func (j *JobMeta) GetLabels() map[string]string {
	return j.Labels
}

func (j *JobMeta) GetRun() RunSpec {
	return j.Run
}

// RunSpec describes when and if a job should be executed.
type RunSpec struct {
	When     string            `json:"when,omitempty"`
	If       string            `json:"if,omitempty"`
	Selector map[string]string `json:"selector,omitempty"`
}
