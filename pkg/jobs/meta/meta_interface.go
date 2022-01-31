package meta

type Job interface {
	GetType() string
	GetName() string
	GetLabels() map[string]string
	GetRun() RunSpec
	Execute() error
}
