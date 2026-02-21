package movements

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

type Middleware func(http.HandlerFunc) http.HandlerFunc

func (h *HttpHandler) SetRoutesTo(r chi.Router) {
	r.Get("/movements", h.query)
}

func (h *HttpHandler) query(w http.ResponseWriter, r *http.Request) {

	movements, err := h.svc.Query(r.Context(), r.URL.Query().Get("rut"))
	if err != nil {
		w.Write([]byte("error"))
		w.WriteHeader(http.StatusInternalServerError)
	}

	response, err := json.Marshal(movements)
	if err != nil {
		w.Write([]byte("error"))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(response)
	w.WriteHeader(http.StatusOK)

}
