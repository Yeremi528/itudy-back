package certificates

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
	r.Post("/certificates", h.issueCertificate)
	r.Get("/certificates", h.getUserCertificates)
}

// POST /certificates
// Body: IssueCertificateRequest
func (h *HttpHandler) issueCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req IssueCertificateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errJSON, status := newError(err)
		w.WriteHeader(status)
		w.Write(errJSON)
		return
	}

	if req.UserID == "" || req.ExamID == "" || req.WorkerName == "" {
		errJSON, status := newError(errBadRequest("user_id, exam_id y worker_name son requeridos"))
		w.WriteHeader(status)
		w.Write(errJSON)
		return
	}

	cert, err := h.svc.IssueCertificate(r.Context(), req)
	if err != nil {
		errJSON, status := newError(err)
		w.WriteHeader(status)
		w.Write(errJSON)
		return
	}

	resp, _ := json.Marshal(cert)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// GET /certificates?userID=xxx
func (h *HttpHandler) getUserCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		errJSON, status := newError(errBadRequest("userID es requerido"))
		w.WriteHeader(status)
		w.Write(errJSON)
		return
	}

	certs, err := h.svc.GetByUserID(r.Context(), userID)
	if err != nil {
		errJSON, status := newError(err)
		w.WriteHeader(status)
		w.Write(errJSON)
		return
	}

	resp, _ := json.Marshal(certs)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
