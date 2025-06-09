package server

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go_course_master/review_manager/internal/app"
	"github.com/go_course_master/review_manager/internal/errs"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Data any `json:"data"`
}

func responseError(l *slog.Logger, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	apiErr := toApiError(err)
	l.Error("http response error", "code", apiErr.Code, "message", apiErr.Message, "orgError", err)

	w.WriteHeader(apiErr.Code)
	json.NewEncoder(w).Encode(&ErrorResponse{
		Code:    apiErr.Code,
		Message: apiErr.Message,
	})
}

func responseSuccess(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&SuccessResponse{
		Data: data,
	})
}

// toApiError - мапппит ошибку в ApiError
func toApiError(err error) *errs.ApiError {
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	switch {
	// Server errors
	case errors.As(err, &syntaxErr):
		return errs.NewApiError(http.StatusBadRequest, "invalid JSON syntax at offset "+strconv.FormatInt(syntaxErr.Offset, 10))
	case errors.As(err, &typeErr):
		return errs.NewApiError(http.StatusBadRequest, "invalid type for field "+typeErr.Field+" at offset "+strconv.FormatInt(typeErr.Offset, 10))
	case errors.Is(err, io.EOF):
		return errs.NewApiError(http.StatusBadRequest, "unexpected end of JSON input")

		// App errors
	case errors.Is(err, app.ErrObjectNotFound):
		return errs.NewApiError(http.StatusNotFound, "object not found")
	case errors.Is(err, app.ErrUserNotFound):
		return errs.NewApiError(http.StatusNotFound, "user not found")
	}

	return errs.NewApiError(http.StatusInternalServerError, "internal server error")
}
