package httpErrors

import (
	"errors"
	tgErrors "github.com/sonyamoonglade/delivery-service/pkg/errors/telegram"
	"github.com/sonyamoonglade/delivery-service/pkg/responder"
	"net/http"
	"strings"
)

const (
	BadRequest                = "Bad request"
	InvalidUsername           = "Validation error. Invalid username"
	InvalidPhoneNumber        = "Validation error. Invalid phone_number"
	InternalServerError       = "Internal server error"
	DeliveryAlreadyExists     = "Delivery already exists"
	DeliveryDoesNotExist      = "Delivery does not exist"
	RunnerAlreadyExists       = "Runner already exists"
	CheckServiceIsUnavailable = "Check service is unavailable"
	TimeoutLimitExceeded      = "Timeout limit exceeded"
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

	msg := strings.ToLower(e.Error())
	switch {

	case strings.Contains(msg, "validation"):
		return BadRequestError(e.Error())
	case strings.Contains(msg, "internal telegram"):
		return InternalTelegramError()
	case strings.Contains(msg, "cli"):
		return NewHttpError(503, CheckServiceIsUnavailable)
	case strings.Contains(msg, "timeout"):
		return NewHttpError(408, TimeoutLimitExceeded)
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
