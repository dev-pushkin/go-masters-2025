package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go_course_master/review_manager/internal/errs"
	"github.com/go_course_master/review_manager/internal/models"
	"github.com/gorilla/mux"
)

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	responseSuccess(w, map[string]string{"status": "ok", "ts": time.Now().Format(time.RFC3339)})
}

func (s *Server) handleReviewAdd(w http.ResponseWriter, r *http.Request) {
	req := models.ReviewAddRequest{}
	// Decode the request body into the ReviewAddRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(s.log, w, err)
		return
	}

	err := s.app.ReviewAdd(r.Context(), req)
	if err != nil {
		responseError(s.log, w, err)
		return
	}
	// If the review was added successfully, respond with a success message
	responseSuccess(w, map[string]string{"message": "Review added successfully"})
}

func (s *Server) handleReviewList(w http.ResponseWriter, r *http.Request) {
	// Decode the query parameters into the ReviewFilter struct
	filter := models.ReviewFilter{}
	// get query parameters
	query := r.URL.Query()
	if userIDStr := query.Get("userID"); userIDStr != "" {
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			responseError(s.log, w, errs.NewApiError(http.StatusBadRequest, "invalid userID format"))
			return
		}
		filter.UserID = &userID
	}
	if objectIDStr := query.Get("objectID"); objectIDStr != "" {
		objectID, err := strconv.ParseUint(objectIDStr, 10, 64)
		if err != nil {
			responseError(s.log, w, errs.NewApiError(http.StatusBadRequest, "invalid objectID format"))
			return
		}
		filter.ObjectID = &objectID
	}

	reviews, err := s.app.ReviewList(r.Context(), filter)
	if err != nil {
		responseError(s.log, w, err)
		return
	}
	responseSuccess(w, reviews)
}

func (s *Server) handleObjectRate(w http.ResponseWriter, r *http.Request) {
	// Handle rating an object
	// This is a placeholder for actual object rating logic.
	vars := mux.Vars(r)
	objectIDparam := vars["objectID"]
	if objectIDparam == "" {
		responseError(s.log, w, errs.NewApiError(http.StatusBadRequest, "objectID is required"))
		return
	}
	objectID, err := strconv.ParseUint(objectIDparam, 10, 64)
	if err != nil {
		responseError(s.log, w, errs.NewApiError(http.StatusBadRequest, "invalid objectID format"))
		return
	}
	rating, err := s.app.ObjectRate(r.Context(), objectID)
	if err != nil {
		responseError(s.log, w, err)
		return
	}
	// If the object was rated successfully, respond with a success message
	responseSuccess(w, rating)
}
