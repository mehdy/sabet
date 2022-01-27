package v1alpha1

// TypeMeta describes an individual job in an configuration
// with strings representing the type of the object and its API schema version.
// Structures that are versioned or persisted should inline TypeMeta.
type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

// ObjectMeta is metadata that all persisted jobs must have.
type ObjectMeta struct {
	Name   string            `json:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
	Run    RunSpec           `json:"run,omitempty"`
}

// RunSpec describes when and if a job should be executed.
type RunSpec struct {
	When     string            `json:"when,omitempty"`
	If       string            `json:"if,omitempty"`
	Selector map[string]string `json:"selector,omitempty"`
}
