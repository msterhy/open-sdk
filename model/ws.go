package model

type Trainer struct {
	Content   string `json:"content"`
	StartTime int64  `json:"start_time"`
	Read      bool   `json:"read"`
	From      string `json:"from"`
}
