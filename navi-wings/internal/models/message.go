	package models

	type Message struct {
		ID     uint   `json:"id" gorm:"primaryKey"`
		Text   string `json:"text"`
		Sender string `json:"sender"`
	}
