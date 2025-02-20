package models

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Attachment  string `json:"attachment,omitempty"`
}