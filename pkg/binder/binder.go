package binder

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"io"
	"reflect"
	"strings"
	"sync"
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

// BindPayload
// Validator for http createDelivery dto (Payload)
func BindPayload(r io.Reader) (*tgdelivery.Payload, error) {

	var p tgdelivery.Payload

	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(bytes, &p); err != nil {
		return nil, err
	}

	if p.Order == nil || p.User == nil {
		return nil, &BindingError{
			Err:     bindingError,
			Message: "order or user fields are missing",
		}
	}

	usrStruct := structs.New(p.User)
	ordStruct := structs.New(p.Order)

	var bindResults []string

	for k, _ := range usrStruct.Map() {
		if val := usrStruct.Field(k); val.IsZero() {
			bindResults = append(bindResults, val.Tag("json"))
		}
	}
	for k, _ := range ordStruct.Map() {
		if val := ordStruct.Field(k); val.IsZero() {
			bindResults = append(bindResults, val.Tag("json"))
		}
	}
	if len(bindResults) > 0 {
		msg := ""
		for i, r := range bindResults {
			splByComma := strings.Split(r, ",")
			if len(splByComma) > 1 {
				r = splByComma[0]
			}
			if i != len(bindResults)-1 {
				msg += fmt.Sprintf("%s, ", strings.ToLower(r))
			} else {
				msg += fmt.Sprintf("%s ", strings.ToLower(r))
			}
		}
		if len(bindResults) == 1 {
			msg += "field is missing in request body"
		} else {
			msg += "fields are missing in request body"
		}
		return nil, &BindingError{
			Err:     bindingError,
			Message: msg,
		}
	}

	return &p, nil
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
	fmt.Println(ptr)
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
