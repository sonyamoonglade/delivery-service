package callback

import (
	"encoding/json"
	"strings"
)

type Data struct {
	DeliveryID int64 `json:"delivery_id,omitempty"`
}

func Decode(data string, out interface{}) error {

	err := json.NewDecoder(strings.NewReader(data)).Decode(&out)
	if err != nil {
		return err
	}
	return nil
}
