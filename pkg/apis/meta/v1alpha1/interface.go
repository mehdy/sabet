package v1alpha1

type Job interface {
	GetName() string
	GetLabels() map[string]string
	GetRun() RunSpec
}
