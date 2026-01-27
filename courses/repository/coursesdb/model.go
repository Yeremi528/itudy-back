package coursesdb

import (
	"github.com/Yeremi528/itudy-back/courses"
)

// --- Estructuras de la Colección 'tracks' ---

// TechAvailability define la estructura de salida: una tecnología y sus niveles disponibles
type TechAvailability struct {
	Tech   string   `bson:"tech" json:"tech"`
	Levels []string `bson:"levels" json:"levels"`
}

// Course (El nodo ligero que está dentro del Track)
type Course struct {
	CourseID    string   `bson:"course_id" json:"course_id"`
	Title       string   `bson:"title" json:"title"`
	Description string   `bson:"description" json:"description"`
	Icon        string   `bson:"icon" json:"icon"`
	SectionID   string   `bson:"section_id" json:"section_id"`
	Order       int      `bson:"order" json:"order"`
	TopicTags   []string `bson:"topic_tags" json:"topic_tags"`
}
type Language struct {
	ID       string `bson:"_id,omitempty"`
	Code     string `bson:"code"`
	Name     string `bson:"name"`
	Flag     string `bson:"flag"`
	IsActive bool   `bson:"is_active"`
}

// Track (El documento principal que contiene el path de cursos)
type Track struct {
	ID        string   `bson:"_id" json:"id"`
	Tech      string   `bson:"tech" json:"tech"`
	Lang      string   `bson:"lang" json:"lang"`
	Level     string   `bson:"level" json:"level"`
	CreatedAt string   `bson:"created_at" json:"created_at"`
	Courses   []Course `bson:"courses" json:"courses"`
}

// --- Estructuras de la Colección 'course_contents' ---

// Exercise (Una pregunta dentro del contenido)
type Exercise struct {
	ID                 string   `bson:"id" json:"id"`
	Type               string   `bson:"type" json:"type"`
	Question           string   `bson:"question" json:"question"`
	CodeSnippet        string   `bson:"code_snippet,omitempty" json:"code_snippet,omitempty"`
	Options            []string `bson:"options" json:"options"`
	CorrectAnswerIndex int      `bson:"correct_answer_index" json:"correct_answer_index"`
	Explanation        string   `bson:"explanation" json:"explanation"`
}

// Content (El documento pesado que contiene teoría y ejercicios)
type Content struct {
	ID             string     `bson:"_id" json:"id"`
	CourseRefID    string     `bson:"course_ref_id" json:"course_ref_id"`
	TheoryMarkdown string     `bson:"theory_markdown" json:"theory_markdown"`
	Exercises      []Exercise `bson:"exercises" json:"exercises"`
}

func trackToCourse(track Track) courses.CourseByID {
	return courses.CourseByID{
		Lang: track.Lang,
		Lv:   track.Level,
		ID:   track.Tech,
		Courses: func() []courses.Course {
			coursesList := make([]courses.Course, len(track.Courses))
			for i, c := range track.Courses {
				coursesList[i] = courses.Course{
					CourseID:    c.CourseID,
					Title:       c.Title,
					Description: c.Description,
					Icon:        c.Icon,
					SectionID:   c.SectionID,
					Order:       c.Order,
					TopicTags:   c.TopicTags,
				}
			}
			return coursesList
		}(),
	}
}

func contentToCourse(content Content) courses.Content {
	return courses.Content{
		ID:          content.ID,
		CourseRefID: content.CourseRefID,
		Theory:      content.TheoryMarkdown,
		Exercises: func() []courses.Exercise {
			exercisesList := make([]courses.Exercise, len(content.Exercises))
			for i, e := range content.Exercises {
				exercisesList[i] = courses.Exercise{
					ID:                 e.ID,
					Type:               e.Type,
					Question:           e.Question,
					Options:            e.Options,
					CorrectAnswerIndex: e.CorrectAnswerIndex,
					Explanation:        e.Explanation,
				}
			}
			return exercisesList
		}(),
	}
}

func languageToLanguage(lang []Language) []courses.Language {

	languages := make([]courses.Language, len(lang))
	for i, l := range lang {
		languages[i] = courseLanguageToLanguage(l)
	}
	return languages
}

func courseLanguageToLanguage(lang Language) courses.Language {
	return courses.Language{
		ID:       lang.ID,
		Code:     lang.Code,
		Name:     lang.Name,
		Flag:     lang.Flag,
		IsActive: lang.IsActive,
	}
}
