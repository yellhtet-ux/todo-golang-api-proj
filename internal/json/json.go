package json

import (
	"encoding/json"
	"net/http"
)

type BaseErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Data    any    `json:"data"`
}

// WRITE
func Write(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func Read(r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func InvalidRequest(w http.ResponseWriter, err error, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(w).Encode(
		BaseErrorResponse{
			Message: "invalid request error",
			Error:   err.Error(),
			Data:    data,
		})
}

func InternalServerError(w http.ResponseWriter, err error, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	return json.NewEncoder(w).Encode(BaseErrorResponse{
		Message: "internal server error",
		Error:   err.Error(),
		Data:    data,
	})
}

func Conflict(w http.ResponseWriter, err error, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	return json.NewEncoder(w).Encode(BaseErrorResponse{
		Message: "conflict error",
		Error:   err.Error(),
		Data:    data,
	})
}

func NotAuthorizedError(w http.ResponseWriter, err error, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	return json.NewEncoder(w).Encode(BaseErrorResponse{
		Message: "missing authorization error",
		Error:   err.Error(),
		Data:    data,
	})
}

func NotFound(w http.ResponseWriter, err error, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	return json.NewEncoder(w).Encode(
		BaseErrorResponse{
			Message: "not found error",
			Error:   err.Error(),
			Data:    data,
		})
}
