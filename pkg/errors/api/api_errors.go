package api_erros

const (
	InternalServerError   = "internal server error"
	DeliveryAlreadyExists = "delivery already exists"
)

type InvalidDelivery struct {
	Err error
}

func (e InvalidDelivery) Error() string {
	return e.Err.Error()
}
