package tgErrors

import "errors"

var (
	InternalError = errors.New("internal telegram error")
)

const (
	RunnerDoesNotExist = "Курьера с таким номером пока не существует, прости 🙇."
)
