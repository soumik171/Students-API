package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	// If wanna use custom struct while viewing in postman, then add struct tags

	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

// For collecting response

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	// For converting struct data to json data, we have to encode

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}
