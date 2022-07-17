package responder

import (
	"encoding/json"
	"net/http"
)

type R map[string]interface{}

func JSON(w http.ResponseWriter, code int, v R) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	bytes, _ := json.Marshal(v)
	w.Write(bytes)

}
