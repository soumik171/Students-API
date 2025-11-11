package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	// must provide the location of folder types, not inbuilt types like:"go/types"
	"github.com/go-playground/validator/v10"
	"github.com/soumik171/Students-API/internal/storage"
	"github.com/soumik171/Students-API/internal/types"
	"github.com/soumik171/Students-API/internal/utils/response"
)

// dependency inject from Storage struct

// Creating new student
func Create(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors) // typecast to slice
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))

			return
		}

		// pass the value to struct
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully,", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		// w.Write([]byte("welcome to students api")) // convert string into byte & pass that to Write()

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})

	}

}

// Student Info Get by Id:

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)

		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")

		students, err := storage.GetStudentsList()

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusOK, students)

	}

}

// Update the student info:

func UpdateById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the id from the url

		id := r.PathValue("id")
		slog.Info("updating a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Decode the request body into Student struct

		var student types.Student

		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// update db:

		updatedStudent, err := storage.UpdateStudentInfo(intId, student)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		// return the updated data as JSON

		response.WriteJson(w, http.StatusOK, updatedStudent)

		slog.Info("student data updated successfully", slog.String("id", id))

	}
}

// Delete the student by id:

func DeleteById(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err = s.DeleteStudent(id)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, fmt.Sprintf("student with id %d deleted successfully", id))

		slog.Info("data deleted succesfully of user", slog.String("id", idStr))
	}
}
