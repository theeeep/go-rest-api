package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/theeeep/go-rest-api/internal/storage"
	"github.com/theeeep/go-rest-api/internal/types"
	"github.com/theeeep/go-rest-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a new student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			slog.Error("No student data provided")
			http.Error(w, "No student data provided", http.StatusBadRequest)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("No student data provided")))
			return
		}

		if err != nil {
			slog.Error("Failed to decode student data", slog.String("error", err.Error()))
			http.Error(w, "Failed to decode student data", http.StatusBadRequest)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Request validation
		if err := validator.New().Struct(student); err != nil {

			validateErrs := err.(validator.ValidationErrors)

			slog.Error("Invalid student data", slog.String("error", err.Error()))

			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// Save student to storage
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("Successfully created a new student", slog.Int64("id", lastId))

		if err != nil {
			slog.Error("Failed to save student", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"id": fmt.Sprintf("Successfully created a new student: %d", lastId)})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting student by id", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentByID(intId)
		if err != nil {
			slog.Error("Failed to get student", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students")
		students, err := storage.GetStudents()
		if err != nil {
			slog.Error("Failed to get students", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, students)
	}
}
