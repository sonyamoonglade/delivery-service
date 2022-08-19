package binder

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/sonyamoonglade/delivery-service/pkg/errors/httpErrors"
)

type BindingError struct {
	Message string
	Err     error
}

var bindingError = errors.New("binding error")

func (e BindingError) Error() string {
	return e.Message
}

var v *validator.Validate

func init() {
	var once sync.Once
	once.Do(func() {
		v = validator.New()
	})
}

type S struct {
	Name string `validate:"required"`
}

func Bind(r io.Reader, out interface{}) error {

	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	//Scan for original type
	typ := reflect.TypeOf(out)
	typDest := reflect.New(typ).Interface()
	if err = json.Unmarshal(bytes, &typDest); err != nil {
		return httpErrors.BadRequestError(httpErrors.BadRequest)
	}
	//Local reflect.Value
	localV := reflect.Indirect(reflect.ValueOf(typDest).Elem())
	ptr := reflect.Indirect(reflect.ValueOf(typDest)).Interface()

	err = v.Struct(ptr)
	if err != nil {
		msg := "Validation error."

		if _, ok := err.(*validator.InvalidValidationError); ok {
			return &BindingError{
				Message: msg,
				Err:     bindingError,
			}
		}

		//todo: snake-case converter
		for _, err := range err.(validator.ValidationErrors) {
			errMsg := err.Error()
			spl := strings.Split(errMsg, "Error:")
			msg += " " + spl[1]
			return &BindingError{
				Message: msg,
				Err:     bindingError,
			}
		}
	}

	//Outer reflect.Value (comes with 'out')
	outV := reflect.Indirect(reflect.ValueOf(out))

	//Change outer to local
	outV.Set(localV)

	return nil
}
