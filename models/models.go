package models

import "time"

type User struct {
	id       uint   `gorm:"primaryKey"`
	login    string `gorm:"unique"`
	password string
}

type Status string

const (
	New        Status = "new"
	InProgress Status = "in progress"
	Done       Status = "done"
)

type Task struct {
	id          uint `gorm:"primaryKey"`
	title       string
	description string
	status      Status `gorm:"type:varchar(11);not null"`
	created_at  time.Time
	updated_at  time.Time
	user_id     uint
}

type Event struct {
	id      uint `gorm:"primaryKey"`
	title   string
	start   time.Time
	end     time.Time
	task_id uint
}
