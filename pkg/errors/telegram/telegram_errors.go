package tgErrors

import "errors"

var (
	InternalError = errors.New("internal telegram error")
)
