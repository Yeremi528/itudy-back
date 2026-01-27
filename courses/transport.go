package courses

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

	r.Get("/courses", h.getAllCourses)
	r.Get("/courses/byID", h.getCourseByID)
	r.Get("/courses/{courseID}/content", h.getCourseContent)
	r.Get("/courses/languages", h.getLanguages)
}

func (h *HttpHandler) getAllCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	courses, err := h.svc.GetAllCourses(r.Context(), r.URL.Query().Get("lang"))
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	coursesJSON, err := json.Marshal(courses)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	w.Write(coursesJSON)
	w.WriteHeader(http.StatusOK)
}

func (h *HttpHandler) getLanguages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	languages, err := h.svc.GetLanguages(r.Context())
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	languagesJSON, err := json.Marshal(languages)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	w.Write(languagesJSON)
	w.WriteHeader(http.StatusOK)
}

func (h *HttpHandler) getCourseByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lv := r.URL.Query().Get("lv")
	lang := r.URL.Query().Get("lang")

	tech := r.URL.Query().Get("tech")
	courses, err := h.svc.GetCourseByID(r.Context(), tech, lv, lang)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	coursesJSON, err := json.Marshal(courses)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	w.Write(coursesJSON)
	w.WriteHeader(http.StatusOK)

}

func (h *HttpHandler) getCourseContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	courseID := chi.URLParam(r, "courseID")
	content, err := h.svc.GetCourseContent(courseID)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	contentJSON, err := json.Marshal(content)
	if err != nil {
		errJSON, status := newError(err)
		w.Write(errJSON)
		w.WriteHeader(status)
		return
	}

	w.Write(contentJSON)
	w.WriteHeader(http.StatusOK)
}
