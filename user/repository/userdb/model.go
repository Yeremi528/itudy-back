package userdb

import "time"

type Stats struct {
	TotalXP       int    `json:"total_xp" bson:"total_xp"`
	StreakDays    int    `json:"streak_days" bson:"streak_days"`
	CurrentLeague string `json:"current_league" bson:"current_league"`
}

type UserCoursesInfo struct {
	Active string           `json:"active" bson:"active"`
	List   []EnrolledCourse `json:"list" bson:"list"`
}

type EnrolledCourse struct {
	ID          string  `json:"id" bson:"_id"`
	Name        string  `json:"name,omitempty" bson:"name,omitempty"`
	Progress    float64 `json:"progress,omitempty" bson:"progress,omitempty"`
	IsCompleted bool    `json:"is_completed,omitempty" bson:"is_completed,omitempty"`
}

type User struct {
	ID             string          `json:"id,omitempty" bson:"_id,omitempty"`
	Email          string          `json:"email" bson:"email"`
	Name           string          `json:"name" bson:"name"`
	Phone          string          `json:"phone,omitempty" bson:"phone,omitempty"`
	Country        string          `json:"country" bson:"country"`
	NativeLanguage string          `json:"native_language" bson:"native_language"`
	CreatedAt      time.Time       `json:"created_at" bson:"created_at"`
	ImageURL       string          `json:"image_url,omitempty" bson:"image_url,omitempty"`
	CoursesInfo    UserCoursesInfo `json:"courses_info,omitempty" bson:"courses_info,omitempty"`
	Stats          Stats           `json:"stats,omitempty" bson:"stats,omitempty"`
}
