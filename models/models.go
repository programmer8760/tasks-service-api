package models

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Login    string `gorm:"unique"`
	Password string
}

type Status string

const (
	TaskNew        Status = "new"
	TaskInProgress Status = "in_progress"
	TaskDone       Status = "done"
)

func (s Status) IsValid() bool {
	switch s {
	case TaskNew, TaskInProgress, TaskDone:
		return true
	}
	return false
}

type Task struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Description string
	Status      Status
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
