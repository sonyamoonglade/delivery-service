package runner

type Storage interface {
	IsRunner() (bool, error)
	Register() (int64, error)
	Ban() (int64, error)
}
