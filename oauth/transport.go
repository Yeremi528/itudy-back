package oauth

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
	r.Post("/oauth/google/login", h.googleLogin)
}

func (h *HttpHandler) googleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := h.svc.GoogleLogin(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		errJSON, status := newError(err)
		w.WriteHeader(status)
		w.Write(errJSON)
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		errJSON, status := newError(err)
		w.WriteHeader(status)
		w.Write(errJSON)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)

}
