package models

import (
	"time"
)

type UserProgress struct {
	UserID       		  string        `json:"user_id"           bson:"userId"` // Unique identifier for the user
	LevelProgress		  LevelProgress `json:"level_progress"    bson:"levelProgress"` // Tracks progress within each level
	TotalTimeSpent 		  int64           `json:"total_time_spent"  bson:"totalTimeSpent"` // Total time spent learning (milliseconds)
	Streak        		  int           `json:"streak"            bson:"streak"` // Consecutive days with at least one lesson
	PointsEarned  		  int           `json:"points_earned"     bson:"pointsEarned"` // Total points earned
	Achievements 		  []string      `json:"achievements"      bson:"achievements"` // List of earned achievements
	CurrentCombo  		  int           `json:"current_combo"     bson:"currentCombo"` // Current combo streak
	HighestCombo  		  int           `json:"highest_combo"     bson:"highestCombo"` // Highest combo achieved
	LastLessonDate 		  time.Time     `json:"last_lesson_date"  bson:"lastLessonDate"` // Date of the last completed lesson
	TotalLevelsCompleted  int			`json:"total_levels_completed" bson:"totalLevelsCompleted"`
  }
  
type LevelProgress 		struct {
	CurrentLevel 	int              	`json:"current_level"      bson:"currentLevel"`     // Current level the user is on
	LevelScores  	map[int]float32     `json:"level_scores"       bson:"levelScores"`      // Scores achieved in each level (key: level number, value: score)
  } 
type User struct {
	UserId           string    			`json:"userId,omitempty" bson:"userId,omitempty"`
	Username         string    			`json:"username,omitempty" bson:"username,omitempty"`
	Email            string    			`json:"email,omitempty" bson:"email,omitempty"`
	Password         string    			`json:"password,omitempty" bson:"password,omitempty"`
	RegistrationDate time.Time 			`json:"registration_date,omitempty" bson:"registration_date,omitempty"`
	Name             string    			`json:"name,omitempty" bson:"name,omitempty"`
	Bio              string    			`json:"bio,omitempty" bson:"bio,omitempty"`
	Location         string    			`json:"location,omitempty" bson:"location,omitempty"`
	ProfilePhoto	 float32				`json:"profile_photo,omitempty" bson:"profilePhoto,omitempty"`
	DoB              time.Time 			`json:"dob,omitempty" bson:"dob,omitempty"` // Use time.Time for date of birth
}

func New(userId, username, email, password string) User {
	return User{
		UserId:           userId,
		Username:         username,
		Email:            email,
		Password:         password,
		RegistrationDate: time.Now(),
	}
}

type UpdateUser struct {
	Name         string    `json:"name,omitempty"`
	Bio          string    `json:"bio,omitempty"`
	Location     string    `json:"location,omitempty"`
	DoB          time.Time `json:"dob,omitempty"` // Optional date of birth
	ProfilePhoto float32       `json:"profile_photo,omitempty"`
}
