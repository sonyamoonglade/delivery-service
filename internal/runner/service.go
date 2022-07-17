package runner

type Service interface {
	IsRunner() (bool, error)
	Register() error
	Ban() error
}
