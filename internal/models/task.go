package models

import "time"

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"userId"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `gorm:"default:false" json:"isCompleted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	
	User        User      `gorm:"foreignKey:UserID" json:"user"`
}


