package learning

import (
	"encoding/json"
	"errors"
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

	r.Post("/lessons", h.newLesson)
	r.Get("/lessons", h.getLessonByID)
	r.Put("/lessons", h.updateLessonByID)

}

func (h *HttpHandler) newLesson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var lesson Lesson
	if err := json.NewDecoder(r.Body).Decode(&lesson); err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	if err := h.svc.CreateLesson(r.Context(), lesson); err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
func (h *HttpHandler) getLessonByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		errJSON, _ := newError(errors.New("missing lesson ID"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	courseID := r.URL.Query().Get("courseID")
	if courseID == "" {
		errJSON, _ := newError(errors.New("missing course ID"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	lesson, err := h.svc.GetLessonByID(r.Context(), userID, courseID)
	if err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	lessonJSON, err := json.Marshal(lesson)
	if err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(lessonJSON)
}

func (h *HttpHandler) updateLessonByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		errJSON, _ := newError(errors.New("missing user ID"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	courseID := r.URL.Query().Get("courseID")
	if courseID == "" {
		errJSON, _ := newError(errors.New("missing course ID"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}
	lessonID := r.URL.Query().Get("lessonID")
	if lessonID == "" {
		errJSON, _ := newError(errors.New("missing lesson ID"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	nextID := r.URL.Query().Get("nextID")
	if nextID == "" {
		errJSON, _ := newError(errors.New("missing next lesson ID"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	if err := h.svc.UpdateLesson(r.Context(), userID, courseID, lessonID, nextID); err != nil {
		errJSON, _ := newError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
