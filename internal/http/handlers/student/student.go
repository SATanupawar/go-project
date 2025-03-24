package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/satyampawar/go-project/internal/types"
	"github.com/satyampawar/go-project/internal/utils/responce"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			responce.WriteJson(w, responce.StatusBadRequest, responce.GeneralError(err))
			return
		}

		// validadte request

		err = validator.New().Struct(student)
		if err != nil {
			responce.WriteJson(w, responce.StatusValidationError, responce.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		slog.Info("student created successfully")
		responce.WriteJson(w, http.StatusCreated, map[string]string{"message": "student created successfully"})
	}
}
