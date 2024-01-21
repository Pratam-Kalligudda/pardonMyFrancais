package models

import "time"

type UserProgress struct {
	Username      string `json:"username" bson:"username"`
	LevelProgress []struct {
		LevelName       string `json:"level_name" bson:"level_name"`
		CurrentScore    int    `json:"current_score" bson:"current_score"`
		CurrentSublevel int    `json:"current_sublevel" bson:"current_sublevel"`
		ScoresHistory   []struct {
			Date  time.Time `json:"date" bson:"date"`
			Score int       `json:"score" bson:"score"`
		} `json:"scores_history" bson:"scores_history"`
	} `json:"level_progress" bson:"level_progress"`
}

type User struct {
	Username         string    `json:"username" bson:"username"`
	Email            string    `json:"email" bson:"email"`
	Password         string    `json:"password" bson:"password"`
	RegistrationDate time.Time `json:"registration_date" bson:"registration_date"`
}

func New(username, email, password string) User {
	return User{
		Username:         username,
		Email:            email,
		Password:         password,
		RegistrationDate: time.Now(),
	}
}