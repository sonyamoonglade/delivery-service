package responder

import (
	"encoding/json"
	"net/http"
)

type R map[string]interface{}

type ValidMap interface {
	map[string]interface{} | map[string]bool | R
}

func JSON[T ValidMap](w http.ResponseWriter, code int, v T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	bytes, _ := json.Marshal(v)

	w.Write(bytes)
}
