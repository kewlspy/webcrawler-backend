package models

import "time"

type URL struct {
	ID              uint      `gorm:"primaryKey"`
	Link            string    `gorm:"not null"`
	Status          string
	HTMLVersion     string
	Title           string
	InternalLinks   int
	ExternalLinks   int
	BrokenLinks     int
	HasLoginForm    bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	BrokenLinkItems []BrokenLink
}

type BrokenLink struct {
	ID     uint   `gorm:"primaryKey"`
	URLID  uint
	Link   string
	Status int
}
