package callback

import (
	"encoding/json"
	"strings"
)

type ReserveData struct {
	DeliveryID int64 `json:"delivery_id,omitempty"`
}

type CompleteData struct {
	DeliveryID int64 `json:"delivery_id,omitempty"`
}

func Decode(data string, out interface{}) error {

	err := json.NewDecoder(strings.NewReader(data)).Decode(&out)
	if err != nil {
		return err
	}
	return nil
}
