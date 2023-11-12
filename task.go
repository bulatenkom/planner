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
	Status      taskStatus
}

func (e Task) GetId() string {
	return e.Id
}

type TaskDto struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	Status      string `json:"status"`
}

func NewTask(title, description string) *Task {

	createdAt := time.Now()

	task := &Task{
		Id:          generateTaskKey(title, createdAt),
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
		Status:      New,
	}
	return task
}

func generateTaskKey(title string, createdAt time.Time) string {
	return fmt.Sprintf("%v-%v", title, createdAt.UnixNano())
}

type taskStatus string

const (
	New    taskStatus = "New"
	Active taskStatus = "Active"
	Halt   taskStatus = "Halt"
	Done   taskStatus = "Done"
)

func getStatusDict() []taskStatus {
	return []taskStatus{New, Active, Halt, Done}
}

func ParseTaskStatus(status string) (taskStatus, error) {
	switch status {
	case "New":
		return New, nil
	case "Active":
		return Active, nil
	case "Halt":
		return Halt, nil
	case "Done":
		return Done, nil
	default:
		return "", fmt.Errorf("unexpected TASK status: %v", status)
	}
}
