package responce

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Responce struct {
	Status int
	Error  string
}

const (
	StatusBadRequest          = 400
	StatusInternalServerError = 500
	StatusValidationError     = 403

	messageError               = "bad request"
	messageInternalServerError = "internal server error"
)

func WriteJson(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)

}

func GeneralError(err error) Responce {

	return Responce{
		Status: StatusBadRequest,
		Error:  messageError,
	}
}

func InternalServerError(err error) Responce {
	return Responce{
		Status: StatusInternalServerError,
		Error:  messageInternalServerError,
	}
}

func ValidationError(errs validator.ValidationErrors) Responce {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("%s is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("%s is not valid", err.Field()))
		}
	}
	return Responce{
		Status: StatusValidationError,
		Error:  strings.Join(errMsgs, ", "),
	}
}

// func Success(message string) Responce {
// 	return Responce{
// 		Status:  StatusOK,
// 		Message: messageSuccess,
// 	}
// }
