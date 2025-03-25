package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/satyampawar/go-project/internal/storage"
	"github.com/satyampawar/go-project/internal/types"
	"github.com/satyampawar/go-project/internal/utils/responce"
	"strconv"
)

func New(storage storage.Storage) http.HandlerFunc {
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

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		slog.Info("student created successfully", "id", lastId)

		if err != nil {
			responce.WriteJson(w, responce.StatusInternalServerError, responce.GeneralError(err))
			return
		}

		slog.Info("student created successfully")
		responce.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		slog.Info("get student by id", slog.String("id", id))


		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			responce.WriteJson(w, responce.StatusBadRequest, responce.GeneralError(err))
			return
		}
		student, err := storage.GetStudentByID(idInt)

		if err != nil {
			responce.WriteJson(w, responce.StatusInternalServerError, responce.GeneralError(err))
			return
		}

		responce.WriteJson(w, http.StatusOK, student)
	}
}
