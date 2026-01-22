package models

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Login    string `gorm:"unique"`
	Password string
}

type Status string

const (
	New        Status = "new"
	InProgress Status = "in progress"
	Done       Status = "done"
)

type Task struct {
	Id          uint `gorm:"primaryKey"`
	Title       string
	Description string
	Status      Status `gorm:"type:varchar(11);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      uint
}

type Event struct {
	ID     uint `gorm:"primaryKey"`
	Title  string
	Start  time.Time
	End    time.Time
	TaskID uint
}
