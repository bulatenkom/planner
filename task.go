package main

import (
	"fmt"
	"time"
)

/*
	[Brief Description]

	User can create/modify/delete/schedule a TASK in backlog.

	User creating TASK can give it a title, description.

	TASK is measured in hours and minutes.

	TASK can be marked as finished/done.

*/

type Task struct {
	Id          string
	Title       string
	Description string
	CreatedAt   time.Time
	Done        bool
}

func (e Task) GetId() string {
	return e.Id
}

type TaskDto struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	Done        bool   `json:"done"`
}

func NewTask(title, description string) *Task {

	createdAt := time.Now()

	task := &Task{
		Id:          generateTaskKey(title, createdAt),
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
		Done:        false,
	}
	return task
}

func generateTaskKey(title string, createdAt time.Time) string {
	return fmt.Sprintf("%v-%v", title, createdAt.UnixNano())
}
