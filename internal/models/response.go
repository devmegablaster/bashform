package models

type Response struct {
	FormID  string   `json:"form_id"`
	Answers []Answer `json:"answers"`
}

type Answer struct {
	QuestionID string `json:"question_id"`
	Value      string `json:"value"`
}
