package main

import (
	"fmt"
)

type SlotStore struct {
	store map[string]*Slot
}

func NewSlotStore() SlotStore {
	ss := SlotStore{}
	ss.store = make(map[string]*Slot, 1024)
	return ss
}

func (ss *SlotStore) FindAll() map[string]*Slot {
	return ss.store
}

func (ss *SlotStore) FindById(id string) (*Slot, error) {
	if slot, ok := ss.FindAll()[id]; ok {
		return slot, nil
	} else {
		return nil, fmt.Errorf("could not find slot with id=%q", id)
	}
}

func (ss *SlotStore) Create(slot *Slot) (*Slot, error) {
	if s, ok := ss.store[slot.Id]; ok {
		return s, fmt.Errorf("couldn't create a SLOT. Store already contains ID=%v", s.Id)
	}
	ss.store[slot.Id] = slot
	return slot, nil
}
