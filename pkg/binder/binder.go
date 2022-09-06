package binder

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

func wrapError(msg string) error {
	return fmt.Errorf("binding error: %s", msg)
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
		return err
	}
	//Local reflect.Value
	localV := reflect.Indirect(reflect.ValueOf(typDest).Elem())
	ptr := reflect.Indirect(reflect.ValueOf(typDest)).Interface()

	err = v.Struct(ptr)
	if err != nil {
		msg := "Validation error."

		if _, ok := err.(*validator.InvalidValidationError); ok {
			return wrapError(msg)
		}

		for _, err := range err.(validator.ValidationErrors) {
			errMsg := err.Error()
			spl := strings.Split(errMsg, "Error:")
			msg += " " + spl[1]
			return wrapError(msg)
		}
	}

	//Outer reflect.Value (comes with 'out')
	outV := reflect.Indirect(reflect.ValueOf(out))

	//Change outer to local
	outV.Set(localV)

	return nil
}
