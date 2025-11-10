package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	// must provide the location of folder types, not inbuilt types like:"go/types"
	"github.com/go-playground/validator/v10"
	"github.com/soumik171/Students-API/internal/types"
	"github.com/soumik171/Students-API/internal/utils/response"
)

func Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// In go, we cannot directly pass the json json data, we have to decode that, then pass the data to struct

		slog.Info("creating a student")

		var student types.Student // creating var of Student struct, so that we can call them

		err := json.NewDecoder(r.Body).Decode(&student) // decode the data, that received from the student

		// check error using error package. we can use normally but we have to check, is it that type of error, that we are expecting

		// errors:package, .Is():method to match our expected error

		// io.EOF: if body is empty, then throw, this kind of error

		if errors.Is(err, io.EOF) {
			// response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))	// For Auto defined Error
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body"))) // For custom error

			return

		}

		// For catching general error:

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		// Request Validation:-->

		// missing element(error) validate:

		errVal := validator.New().Struct(student)
 
		if errVal != nil {
			validateErrs := errVal.(validator.ValidationErrors) // typecast to slice
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))

			return
		}

		// w.Write([]byte("welcome to students api")) // convert string into byte & pass that to Write()

		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "OK"})

	}
}
