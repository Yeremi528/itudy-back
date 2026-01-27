package courses

import (
	"encoding/json"
	"net/http"

	"github.com/Yeremi528/itudy-back/kit/web"
)

// TechAvailability define la estructura de salida: una tecnología y sus niveles disponibles
type TechAvailability struct {
	Tech   string   `bson:"tech" json:"tech"`
	Levels []string `bson:"levels" json:"levels"`
}

type CourseByID struct {
	Lang    string   `json:"lang"`
	Lv      string   `json:"lv"`
	ID      string   `json:"id"`
	Courses []Course `json:"courses"`
}
type Language struct {
	ID       string `json:"_id,omitempty"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Flag     string `json:"flag"`
	IsActive bool   `json:"isActive"`
}

type Course struct {
	CourseID    string   `json:"course_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	SectionID   string   `json:"section_id"`
	Order       int      `json:"order"`
	Status      string   `json:"status"` // "locked", "active", "completed"
	TopicTags   []string `json:"topic_tags"`
}

type Content struct {
	ID          string     `json:"_id"`
	CourseRefID string     `json:"course_ref_id"`
	Theory      string     `json:"theory_markdown"`
	Exercises   []Exercise `json:"exercises"`
}

type Exercise struct {
	ID                 string   `json:"id"`
	Type               string   `json:"type"`
	Question           string   `json:"question"`
	Options            []string `json:"options"`
	CorrectAnswerIndex int      `json:"correct_answer_index"`
	Explanation        string   `json:"explanation"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}

func newError(err error) ([]byte, int) {
	var status int
	switch {
	case web.IsRequestError(err):
		errReq := web.GetRequestError(err)
		status = errReq.Status
	default:
		status = http.StatusInternalServerError
	}

	errorResponse := ErrorResponse{
		ErrorMessage: err.Error(),
	}
	errorResponseJSON, err := json.Marshal(errorResponse)
	if err != nil {
		return []byte(err.Error()), http.StatusInternalServerError
	}

	return errorResponseJSON, status
}
