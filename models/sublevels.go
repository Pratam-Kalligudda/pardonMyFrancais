package models

type Sublevel struct {
	LevelName string     `json:"level_name" bson:"level_name"`
	Sublevels []Question `json:"questions" bson:"questions"`
	Sentence  string     `json:"sentence" bson:"sentence"`
}

type Question struct {
	Question      string   `json:"question" bson:"question"`
	Options       []string `json:"options" bson:"options"`
	CorrectOption string   `json:"correct_option" bson:"correct_option"`
}
