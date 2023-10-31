package main

import (
	"fmt"
	"time"
)

/*
	[Brief Description]

	Planner is an webapp that provides an API and UI that allows user to schedule his dayly routines and track its progress

	[Actors]

	User can allocate time EVENT in any range, but most likely within a day (at this point EVENT can be treated as a routine).

	User can track its routines using different representations month/week/day.

	User creating EVENT can give it a title, description and duration.

	EVENT is measured in hours and minutes.

	EVENT can be of type:
		- Daily routine
		- Education
		- Regular Event
		- One-time Event

	EVENT can be marked as finished/done.

	EVENT must be assigned to date and time.

*/

type Event struct {
	Id          string
	Title       string
	Description string
	Duration    time.Duration
	PlannedOn   time.Time
	CreatedAt   time.Time
	Type        eventType
	Done        bool
}

func (e Event) GetId() string {
	return e.Id
}

type EventDto struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    string `json:"duration"`
	PlannedOn   string `json:"plannedOn"`
	CreatedAt   string `json:"createdAt"`
	Type        string `json:"type"`
	Done        bool   `json:"done"`
}

func NewEvent(
	title, description string,
	plannedOn time.Time,
	duration time.Duration,
	eventType eventType) *Event {

	createdAt := time.Now()

	event := &Event{
		Id:          generateEventKey(eventType, title, createdAt),
		Title:       title,
		Description: description,
		Duration:    duration,
		PlannedOn:   plannedOn,
		CreatedAt:   createdAt,
		Type:        eventType,
		Done:        false,
	}
	return event
}

func generateEventKey(eventType eventType, title string, createdAt time.Time) string {
	return fmt.Sprintf("%v-%v-%v", eventType, title, createdAt.UnixNano())
}

type eventType string

const (
	OneTimeEvent eventType = "OneTimeEvent"
	RegularEvent eventType = "RegularEvent"
	DailyRoutine eventType = "DailyRoutine"
	Education    eventType = "Education"
)

func ParseEventType(eventType string) (eventType, error) {
	switch eventType {
	case "OneTimeEvent":
		return OneTimeEvent, nil
	case "RegularEvent":
		return RegularEvent, nil
	case "DailyRoutine":
		return DailyRoutine, nil
	case "Education":
		return Education, nil
	default:
		return "", fmt.Errorf("unexpected EVENT type: %v", eventType)
	}
}
