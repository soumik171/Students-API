package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
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

// For general error:
func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

// For validate error:
func ValidationError(errs validator.ValidationErrors) Response {
	// return the list of errors
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required feild", err.Field())) //return the fieldname of the struct

		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid field", err.Field()))

		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ","), // By using string packages, can convert the slice elements into strings
	}
}
