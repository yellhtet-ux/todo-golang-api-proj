package json

import (
	"encoding/json"
	"net/http"
)

// WRITE
func Write (w http.ResponseWriter,status int, data any) error {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func Read (r *http.Request,data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func InvalidRequest(w http.ResponseWriter,err error) error {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(w).Encode(map[string]string{"message":"Invalid request","error":err.Error()})
}

func InternalServerError(w http.ResponseWriter,err error) error {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusInternalServerError)
	return json.NewEncoder(w).Encode(map[string]string{"message":"Internal server error","error":err.Error()})
}

func NotFound(w http.ResponseWriter,err error) error {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusNotFound)
	return json.NewEncoder(w).Encode(map[string]string{"message":"Not found","error":err.Error()})
}
