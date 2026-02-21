package payin

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

	r.Post("/payin", h.payin)

}

func (h *HttpHandler) payin(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("data.id")
	topic := r.URL.Query().Get("type")

	if topic == "" || ID == "" {
		var body MPWebhookBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			errJSON, status := newError(err)
			w.Write(errJSON)
			w.WriteHeader(status)
			return
		}
		topic = body.Type
		ID = body.Data.ID
	}

	w.Header().Set("Content-Type", "application/json")
	if err := h.svc.WebHook(r.Context(), ID, topic); err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return

	}

	w.WriteHeader(http.StatusOK)
}
