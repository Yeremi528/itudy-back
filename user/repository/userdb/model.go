package userdb

import "time"

type Stats struct {
	TotalXP       int    `json:"total_xp" bson:"total_xp"`
	StreakDays    int    `json:"streak_days" bson:"streak_days"`
	CurrentLeague string `json:"current_league" bson:"current_league"`
}

type User struct {
	ID             string    `json:"id,omitempty" bson:"_id,omitempty"`
	Email          string    `json:"email" bson:"email"`
	Name           string    `json:"name" bson:"name"`
	Phone          string    `json:"phone,omitempty" bson:"phone,omitempty"`
	Country        string    `json:"country" bson:"country"`
	NativeLanguage string    `json:"native_language" bson:"native_language"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	ImageURL       string    `json:"image_url,omitempty" bson:"image_url,omitempty"`
	Stats          *Stats    `json:"stats,omitempty" bson:"stats,omitempty"`
}
