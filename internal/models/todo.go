package models

type Todo struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Attachment  string `json:"attachment,omitempty"`
}