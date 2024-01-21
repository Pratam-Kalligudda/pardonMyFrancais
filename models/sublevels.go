package models

type Sublevel struct {
	LevelName      string   `json:"level_name" bson:"level_name"`
	Sublevels      []Subitem `json:"sublevels" bson:"sublevels"`
}

type Subitem struct {
	Question       string   `json:"question" bson:"question"`
	Options        []string `json:"options" bson:"options"`
	CorrectOption  string   `json:"correct_option" bson:"correct_option"`
}
