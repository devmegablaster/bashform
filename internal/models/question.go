package models

type Question struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Type    string   `json:"type"`
	Options []Option `json:"options"`
}

type QuestionRequest struct {
	Text    string          `json:"text"`
	Type    string          `json:"type"`
	Options []OptionRequest `json:"options"`
}
