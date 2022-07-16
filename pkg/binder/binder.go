package binder

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	tgdelivery "github.com/sonyamoonglade/delivery-service"
	"io"
	"strings"
)

var (
	//bindingError shows that some value is missing
	bindingError = "invalid request body"
)

type BindingError struct {
	Message string
	Err     error
}

func (e *BindingError) Unwrap() error {
	return e.Err
}

func (e *BindingError) Error() string {
	return e.Message
}

// Bind
// Validator for http createDelivery dto (Payload)
func Bind(r io.Reader) (*tgdelivery.Payload, error) {

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
			Message: "Order or User fields are missing",
			Err:     errors.New(bindingError),
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
			Message: msg,
			Err:     errors.New(bindingError),
		}
	}

	return &p, nil
}
