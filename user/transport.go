package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	r.Put("/user", h.updateUser)
	r.Put("/user/image", h.uploadProfileImage)
	r.Post("/user/achievement", h.addAchievement)
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

func (h *HttpHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}
	if err := h.svc.UpdateUser(r.Context(), user); err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *HttpHandler) uploadProfileImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseMultipartForm(10 << 20); err != nil { // límite 10 MB
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		userID = r.FormValue("user_id")
	}
	if userID == "" {
		errJSON, _ := newError(errBadRequest("userID requerido"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	url, err := h.svc.UploadProfileImage(r.Context(), userID, header.Filename, file, contentType)
	if err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	resp := struct {
		ImageURL string `json:"image_url"`
	}{ImageURL: url}
	respJSON, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)
}

func (h *HttpHandler) addAchievement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body struct {
		UserID   string `json:"user_id"`
		Title    string `json:"title"`
		ExamName string `json:"exam_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	achievement := Achievement{
		ID:       uuid.New().String(),
		Title:    body.Title,
		ExamName: body.ExamName,
		EarnedAt: time.Now().UTC(),
	}

	if err := h.svc.AddAchievement(r.Context(), body.UserID, achievement); err != nil {
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
