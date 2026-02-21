package assignments

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HttpHandler struct {
	svc Service
}

func MakeHandlerWith(svc Service) *HttpHandler {
	return &HttpHandler{svc: svc}
}

func (h *HttpHandler) SetRoutesTo(r chi.Router) {
	r.Get("/assignments/availables", h.assignmentsAvailables)
	r.Get("/assignments", h.taskByUserID)
	r.Post("/assignments", h.CreateAssignmentsByUserID)

}

func (h *HttpHandler) assignmentsAvailables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tech := r.URL.Query().Get("tech")
	level := r.URL.Query().Get("level")

	assignmentsAvailables, err := h.svc.AssignmentsAvailables(r.Context(), tech, level)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	assignmentsAvailablesJSON, err := json.Marshal(assignmentsAvailables)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	w.Write(assignmentsAvailablesJSON)
	w.WriteHeader(http.StatusOK)
}

func (h *HttpHandler) taskByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.URL.Query().Get("userID")
	AssignmentTestByUserID, err := h.svc.AssignmentTestByUserID(r.Context(), userID)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	AssignmentTestByUserIDJSON, err := json.Marshal(AssignmentTestByUserID)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	w.Write(AssignmentTestByUserIDJSON)
	w.WriteHeader(http.StatusOK)

}

func (h *HttpHandler) CreateAssignmentsByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newAssignment AssignmentTest
	err := json.NewDecoder(r.Body).Decode(&newAssignment)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	linkPago, err := h.svc.CreateAssignmentsByUserID(r.Context(), newAssignment)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	bytes, err := json.Marshal(responseMobile(linkPago))
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	w.Write(bytes)
	w.WriteHeader(http.StatusCreated)
}
