package api

import (
	"crypto/aes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go_course_master/homework/hw_01/internal/errs"
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
	l.Error("http response error", "code", apiErr.Code, "message", apiErr.Message)

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
	switch e := err.(type) {
	case *errs.ApiError:
		return e
	case *json.UnmarshalTypeError:
		return errs.NewApiError(http.StatusBadRequest, "invalid type")
	case aes.KeySizeError:
		return errs.NewApiError(http.StatusBadRequest, "invalid key size")
	}

	if errors.Is(err, io.ErrUnexpectedEOF) {
		return errs.NewApiError(http.StatusBadRequest, "empty data")
	}
	if errors.Is(err, io.EOF) {
		return errs.NewApiError(http.StatusBadRequest, "empty data")
	}
	if errors.Is(err, io.ErrShortBuffer) {
		return errs.NewApiError(http.StatusBadRequest, "buffer too short")
	}

	return errs.NewApiError(http.StatusInternalServerError, "internal server error")
}
