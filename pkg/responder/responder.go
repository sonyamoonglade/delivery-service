package responder

import (
	"encoding/json"
	"net/http"
)

type R map[string]string

func JSON(r http.ResponseWriter, v R, code int) {
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(code)
	bytes, _ := json.Marshal(v)
	r.Write(bytes)

}
