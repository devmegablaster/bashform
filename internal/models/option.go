package models

type Option struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type OptionRequest struct {
	Text string `json:"text"`
}
