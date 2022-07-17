package httpErrors

import (
	"errors"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/responder"
	"net/http"
	"strings"
)

const (
	InternalServerError   = "Internal server error"
	DeliveryAlreadyExists = "Delivery already exists"
	DeliveryDoesNotExist  = "Delivery does not exist"
)

type HttpError struct {
	code int
	err  string
}

func (e HttpError) Error() string {
	return e.err
}

func (e HttpError) Code() int {
	return e.code
}

func NewHttpError(code int, e string) HttpError {
	return HttpError{
		code: code,
		err:  e,
	}
}

func ConflictError(e string) HttpError {
	return HttpError{
		code: http.StatusConflict,
		err:  e,
	}
}

func ForbiddenError(e string) HttpError {
	return HttpError{
		code: http.StatusForbidden,
		err:  e,
	}
}

func BadRequestError(e string) HttpError {
	return HttpError{
		code: http.StatusBadRequest,
		err:  e,
	}
}

func NotFoundError(e string) HttpError {
	return HttpError{
		code: http.StatusNotFound,
		err:  e,
	}
}

func InternalError() HttpError {
	return HttpError{
		code: http.StatusInternalServerError,
		err:  InternalServerError,
	}
}

func InternalTelegramError() HttpError {
	return HttpError{
		code: http.StatusServiceUnavailable,
		err:  tgErrors.InternalError.Error(),
	}
}

func parseError(e error) HttpError {

	switch {

	case strings.Contains(strings.ToLower(e.Error()), "missing"):
		return BadRequestError(e.Error())
	case strings.Contains(strings.ToLower(e.Error()), "internal telegram"):
		return InternalTelegramError()

	default:
		return InternalError()
	}

}

func baseErrResponse(m string, code int) responder.R {
	return responder.R{
		"message":    m,
		"statusCode": code,
	}
}

func Response(e error) (int, responder.R) {
	var httpErr HttpError
	if errors.As(e, &httpErr) {
		return httpErr.code, baseErrResponse(httpErr.Error(), httpErr.Code())
	}

	parErr := parseError(e)

	return parErr.code, baseErrResponse(parErr.Error(), parErr.Code())

}
