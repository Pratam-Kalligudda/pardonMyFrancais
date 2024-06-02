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
	Username         string  `json:"username,omitempty" bson:"username,omitempty"`
	Email            string  `json:"email,omitempty" bson:"email,omitempty"`
	Password         string  `json:"password,omitempty" bson:"password,omitempty"`
	RegistrationDate time.Time `json:"registration_date,omitempty" bson:"registration_date,omitempty"`
	Name             string  `json:"name,omitempty" bson:"name,omitempty"`
	Bio              string  `json:"bio,omitempty" bson:"bio,omitempty"`
	Location         string  `json:"location,omitempty" bson:"location,omitempty"`
	DoB              time.Time `json:"dob,omitempty" bson:"dob,omitempty"` // Use time.Time for date of birth
  }
  
func New(username, email, password string) User {
	return User{
		Username:         username,
		Email:            email,
		Password:         password,
		RegistrationDate: time.Now(),
	}
}

type UpdateUser struct {
	Name  string `json:"name,omitempty"`
	Bio   string `json:"bio,omitempty"`
	Location  string `json:"location,omitempty"`
	DoB    time.Time `json:"dob,omitempty"` // Optional date of birth
  }