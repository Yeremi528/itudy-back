package exam

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HttpHandler struct {
	svc Service
}

func NewHttpHandler(svc Service) *HttpHandler {
	return &HttpHandler{svc: svc}
}

func (h *HttpHandler) SetRoutesTo(r chi.Router) {
	r.Get("/exam", h.examByUserID)

}

func (h *HttpHandler) examByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	email := r.URL.Query().Get("email")
	examByUserID, err := h.svc.ExamByUserID(r.Context(), email)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	examByUserIDJSON, err := json.Marshal(examByUserID)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	w.Write(examByUserIDJSON)
	w.WriteHeader(http.StatusOK)
}
