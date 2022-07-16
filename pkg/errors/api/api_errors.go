package api_erros

var (
	// InternalServerError describes message if something has gone wrong.
	InternalServerError = "Internal server error"
)

type InvalidDelivery struct {
	Err error
}

func (e InvalidDelivery) Error() string {
	return e.Err.Error()
}
