package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/theeeep/go-rest-api/internal/types"
	"github.com/theeeep/go-rest-api/internal/utils/response"
)

func New() http.HandlerFunc {
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

		response.WriteJson(w, http.StatusCreated, map[string]string{"message": "Successfully created a new student"})
	}
}
