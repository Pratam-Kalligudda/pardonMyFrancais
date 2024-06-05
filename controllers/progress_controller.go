package controllers

import (
	"net/http"
	"strconv"
	"time"

	"example.com/backend/configs"
	"example.com/backend/models"
	"github.com/labstack/echo/v4"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type ProgressResponse struct {
    CurrentLevel   int                 `json:"current_level"`
    LevelScores    map[int]float32         `json:"level_scores"`
    AllQuestionsRight bool             `json:"allQuestionsRight"`
    TotalTimeSpent string              `json:"totalTimeSpent"`
}

type Achievement struct {
    Name        string `bson:"name"`
    Description string `bson:"description"`
}
// UpdateUserProgress handles the update of user progress when they complete a level
func UpdateUserProgress(c echo.Context) error {
	userId := c.Get("userId").(string)
	response := new(ProgressResponse)
	if err := c.Bind(response); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	userProgress, err := getUserProgress(c, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	println(response.AllQuestionsRight)
	println(response.TotalTimeSpent)
	updateLevelProgress(userProgress, response)
	updateComboStreak(userProgress, response.AllQuestionsRight)
	updateStreak(userProgress)
	updatePointsEarned(userProgress, response.LevelScores[response.CurrentLevel])
	updateTotalTimeSpent(userProgress, response.TotalTimeSpent)
	updateLastLessonDate(userProgress)
	updateTotalLevelsCompleted(userProgress)
	err = updateAchievements(c,userProgress)
	if err != nil {
		return err
	}

	err = saveUserProgress(c, userProgress)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user progress"})
	}

	return c.JSON(http.StatusOK, userProgress)
}

// getUserProgress fetches the user progress from the database
func getUserProgress(c echo.Context, userID string) (*models.UserProgress, error) {
	progressCollection := configs.GetClient().Database("pardon_my_francais").Collection("user_progress")
	filter := bson.M{"userId": userID}

	var userProgress models.UserProgress
	err := progressCollection.FindOne(c.Request().Context(), filter).Decode(&userProgress)
	if err != nil {
		return nil, err
	}

	return &userProgress, nil
}

// updateLevelProgress updates the level progress of the user
func updateLevelProgress(userProgress *models.UserProgress, level *ProgressResponse) {
	userProgress.LevelProgress.CurrentLevel = level.CurrentLevel
	if userProgress.LevelProgress.LevelScores == nil {
		userProgress.LevelProgress.LevelScores = make(map[int]float32)
	}
	if level.LevelScores[level.CurrentLevel] > userProgress.LevelProgress.LevelScores[level.CurrentLevel] {
		userProgress.LevelProgress.LevelScores[level.CurrentLevel] = level.LevelScores[level.CurrentLevel]
	}
}

// updateComboStreak updates the combo streak of the user
func updateComboStreak(userProgress *models.UserProgress, allQuestionsRight bool) {
	if allQuestionsRight {
		userProgress.CurrentCombo++
		if userProgress.CurrentCombo > userProgress.HighestCombo {
			userProgress.HighestCombo = userProgress.CurrentCombo
		}
	} else {
		userProgress.CurrentCombo = 0
	}
}

// updateStreak updates the streak of the user
func updateStreak(userProgress *models.UserProgress) {
	if time.Since(userProgress.LastLessonDate).Hours() >= 24 {
		if time.Since(userProgress.LastLessonDate).Hours() < 48 {
			userProgress.Streak++
		} else {
			userProgress.Streak = 1 // Reset streak if more than 1 day missed
		}
	}
}

// updatePointsEarned updates the points earned by the user
// updatePointsEarned updates the points earned by the user
func calculatePoints(levelScore float32,currentCombo int, streak int) int {
    // Base points for completing a level
    basePoints := 1000

	// Points based on the user's score for the level
    scoreBasedPoints := levelScore / 2 

    // Combo bonus (e.g., 5 points per combo streak)
    comboBonus := currentCombo * 5

    // Streak bonus (e.g., 20 points per day streak)
    streakBonus := streak * 20

    // Calculate total points
    totalPoints := basePoints + comboBonus + streakBonus + int(scoreBasedPoints)

    return totalPoints
}

func updatePointsEarned(userProgress *models.UserProgress,levelScore float32) {
    points := calculatePoints(levelScore,userProgress.CurrentCombo, userProgress.Streak)
    userProgress.PointsEarned += points
}

// updateTotalTimeSpent updates the total time spent by the user
func updateTotalTimeSpent(userProgress *models.UserProgress, totalTimeSpentStr string) {
	if totalTimeSpentStr != "" {
		totalTimeSpent, err := strconv.ParseInt(totalTimeSpentStr, 10, 64)
		if err == nil {
			userProgress.TotalTimeSpent += totalTimeSpent
		}
	}
}

// updateLastLessonDate updates the last lesson date of the user
func updateLastLessonDate(userProgress *models.UserProgress) {
	userProgress.LastLessonDate = time.Now()
}

// saveUserProgress saves the updated user progress to the database
func saveUserProgress(c echo.Context, userProgress *models.UserProgress) error {
	progressCollection := configs.GetClient().Database("pardon_my_francais").Collection("user_progress")
	filter := bson.M{"userId": userProgress.UserID}
	update := bson.M{
		"$set": userProgress,
	}

	_, err := progressCollection.UpdateOne(c.Request().Context(), filter, update, options.Update().SetUpsert(true))
	return err
}

func fetchAchievements(c echo.Context) ([]Achievement, error) {
    achievementCollection := configs.GetClient().Database("pardon_my_francais").Collection("achievements")
    var achievements []Achievement

    cursor, err := achievementCollection.Find(c.Request().Context(), bson.M{})
    if err != nil {
        return nil, err
    }

    if err := cursor.All(c.Request().Context(), &achievements); err != nil {
        return nil, err
    }
	println(len(achievements))
    return achievements, nil
}

// updateAchievements updates the user's achievements based on their progress
func updateAchievements(c echo.Context, userProgress *models.UserProgress) error {
    achievements, err := fetchAchievements(c)
    if err != nil {
        return err
    }

    achievementMap := map[string]bool{}
    for _, ach := range userProgress.Achievements {
        achievementMap[ach] = true
    }

    for _, ach := range achievements {
        switch ach.Name {
        case "Novice Learner":
            if userProgress.TotalLevelsCompleted >= 1 && !achievementMap[ach.Name] {
                userProgress.Achievements = append(userProgress.Achievements, ach.Name)
            }
        case "Dedicated Learner":
            if userProgress.TotalLevelsCompleted >= 5 && !achievementMap[ach.Name] {
                userProgress.Achievements = append(userProgress.Achievements, ach.Name)
            }
        case "Master Learner":
            if userProgress.TotalLevelsCompleted >= 10 && !achievementMap[ach.Name] {
                userProgress.Achievements = append(userProgress.Achievements, ach.Name)
            }
        case "Streak Beginner":
            if userProgress.Streak >= 1 && !achievementMap[ach.Name] {
                userProgress.Achievements = append(userProgress.Achievements, ach.Name)
            }
        case "Streak Pro":
            if userProgress.Streak >= 7 && !achievementMap[ach.Name] {
                userProgress.Achievements = append(userProgress.Achievements, ach.Name)
            }
        case "Streak Champion":
            if userProgress.Streak >= 30 && !achievementMap[ach.Name] {
                userProgress.Achievements = append(userProgress.Achievements, ach.Name)
            }
        case "Point Collector":
            if userProgress.PointsEarned >= 10000 && !achievementMap[ach.Name] {
                userProgress.Achievements = append(userProgress.Achievements, ach.Name)
            }
        }
    }

    return nil
}


func updateTotalLevelsCompleted(userProgress *models.UserProgress) {
    userProgress.TotalLevelsCompleted = len(userProgress.LevelProgress.LevelScores)
}