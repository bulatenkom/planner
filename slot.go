package main

import (
	"fmt"
	"time"
)

/*
	[Brief Description]

	Planner is an webapp that provides an API and UI that allows user to schedule his dayly routines and track its progress

	[Actors]

	User can allocate time SLOT in any range, but most likely within a day (at this point SLOT can be treated as a routine).

	User can track its routines using different representations month/week/day.

	User creating SLOT can give it a title, description and duration.

	SLOT is measured in hours and minutes.

	SLOT can be of type:
		- Daily routine
		- Education
		- Regular Event
		- One-time Event

	SLOT can be marked as finished/done.

	SLOT must be assigned to date and time.

*/

// ------------------------------------------------------------
// DOMAIN
// ------------------------------------------------------------

type ProfileSettings struct {
	// settings
}

type SlotType string

const (
	OneTimeEvent SlotType = "OneTimeEvent"
	RegularEvent SlotType = "RegularEvent"
	DailyRoutine SlotType = "DailyRoutine"
	Education    SlotType = "Education"
)

func ParseSlotType(slotType string) (SlotType, error) {
	switch slotType {
	case "OneTimeEvent":
		return OneTimeEvent, nil
	case "RegularEvent":
		return RegularEvent, nil
	case "DailyRoutine":
		return DailyRoutine, nil
	case "Education":
		return Education, nil
	default:
		return "", fmt.Errorf("unexpected SLOT type: %v", slotType)
	}
}

func NewSlot(
	title, description string,
	plannedOn time.Time,
	duration time.Duration,
	slotType SlotType) *Slot {

	createdAt := time.Now()

	slot := &Slot{
		Id:          generateKey(slotType, title, createdAt),
		Title:       title,
		Description: description,
		Duration:    duration,
		PlannedOn:   plannedOn,
		CreatedAt:   createdAt,
		Type:        slotType,
		Done:        false,
	}
	return slot
}

func generateKey(slotType SlotType, title string, createdAt time.Time) string {
	return fmt.Sprintf("%v-%v-%v", slotType, title, createdAt.UnixNano())
}

type Slot struct {
	Id          string
	Title       string
	Description string
	Duration    time.Duration
	PlannedOn   time.Time
	CreatedAt   time.Time
	Type        SlotType
	Done        bool
}

type SlotDto struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    string `json:"duration"`
	PlannedOn   string `json:"plannedOn"`
	CreatedAt   string `json:"createdAt"`
	Type        string `json:"type"`
	Done        bool   `json:"done"`
}
