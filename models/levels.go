package models

type GuidebookContent struct {
	LevelName        string `json:"level_name" bson:"level_name"`
	GuidebookContent []struct {
		FrenchWord          string `json:"french_word" bson:"french_word"`
		FrenchPronunciation string `json:"french_pronunciation" bson:"french_pronunciation"`
		EnglishTranslation  string `json:"english_translation" bson:"english_translation"`
	} `json:"guidebook_content" bson:"guidebook_content"`
}
