package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// A handler is similar to a Controller.
type HttpHandler struct {
	svc Service
}

func MakeHandlerWith(svc Service) *HttpHandler {
	return &HttpHandler{svc: svc}
}

func (h *HttpHandler) SetRoutesTo(r chi.Router) {

	r.Post("/user", h.newUser)
	r.Get("/user", h.queryUser)

}

func (h *HttpHandler) newUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}
	if err := h.svc.CreateUser(r.Context(), user); err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (h *HttpHandler) queryUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	email := r.URL.Query().Get("email")
	user, err := h.svc.GetUser(r.Context(), email)
	if err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}
