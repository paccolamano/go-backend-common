package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

const (
	ServiceNilError = "one or more of the services is null"

	HttpRequestInfoMessage  = "sending http %s request to %s with body %v"
	HttpResponseInfoMessage = "received http reply with status code %d and body %v"
	HttpErrorMessage        = "error while sending http %s request to %s caused by %s"

	EntityCreateError   = "cannot create %s caused by %s"
	EntityGetError      = "cannot get %s caused by %s"
	EntityUpdateError   = "cannot update %s caused by %s"
	EntityDeleteError   = "cannot delete %s caused by %s"
	EntityFindError     = "cannot find %s caused by %s"
	EntityNotFoundError = "%s with uuid %s not found"
)

type RestError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func NewRestError(code int, message string, values ...interface{}) *RestError {
	return &RestError{
		Message: fmt.Sprintf(message, values...),
		Code:    code,
	}
}

func NewBadRequestError(message string, values ...interface{}) *RestError {
	return &RestError{
		Message: fmt.Sprintf(message, values...),
		Code:    http.StatusBadRequest,
	}
}

func NewUnauthorizedError(message string, values ...interface{}) *RestError {
	return &RestError{
		Message: fmt.Sprintf(message, values...),
		Code:    http.StatusUnauthorized,
	}
}

func NewNotFoundError(message string, values ...interface{}) *RestError {
	return &RestError{
		Message: fmt.Sprintf(message, values...),
		Code:    http.StatusNotFound,
	}
}

func NewInternalServerError(message string, values ...interface{}) *RestError {
	return &RestError{
		Message: fmt.Sprintf(message, values...),
		Code:    http.StatusInternalServerError,
	}
}

type validationError struct {
	err error
}

func NewValidationError(err error) validationError {
	return validationError{
		err: err,
	}
}

func (q validationError) String() string {
	var sb strings.Builder

	switch typedError := any(q.err).(type) {
	case validator.ValidationErrors:
		sb.WriteString(parseValidationError(typedError))
	case *json.UnmarshalTypeError:
		sb.WriteString(parseMarshallingError(*typedError))
	default:
		if q.err.Error() == "EOF" {
			sb.WriteString("empty request body")
			break
		}
		sb.WriteString(q.err.Error())
	}

	return strings.Trim(sb.String(), " ")
}

func parseValidationError(typedError validator.ValidationErrors) string {
	var sb strings.Builder
	for _, e := range typedError {
		sb.WriteString(fmt.Sprintf("validation failed on field '%s', condition: %s", e.Field(), e.ActualTag()))

		if e.Param() != "" {
			sb.WriteString(" { " + e.Param() + " }")
		}

		if e.Value() != nil && e.Value() != "" {
			sb.WriteString(fmt.Sprintf(", actual: %v", e.Value()))
		}
		sb.WriteString(". ")
	}

	return sb.String()
}

func parseMarshallingError(e json.UnmarshalTypeError) string {
	return fmt.Sprintf("The field %s must be a %s", e.Field, e.Type.String())
}
