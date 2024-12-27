package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/theeeep/go-rest-api/internal/types"
	"github.com/theeeep/go-rest-api/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			slog.Error("No student data provided")
			http.Error(w, "No student data provided", http.StatusBadRequest)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		slog.Info("Creating a new student")

		response.WriteJson(w, http.StatusCreated, map[string]string{"message": "Successfully created a new student"})
	}
}
