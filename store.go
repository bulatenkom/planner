package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// L1 Cache
type SlotStore struct {
	store map[string]*Slot
}

func NewSlotStore() SlotStore {
	ss := SlotStore{}
	ss.store = make(map[string]*Slot, 1024)
	// ------------------------------------------------------------
	// load file into store
	// ------------------------------------------------------------
	dirs, err := os.ReadDir("data/slots")
	if err != nil {
		panic(err)
	}
	for _, entry := range dirs {
		if entry.IsDir() {
			continue
		}
		filedata, err := os.ReadFile("data/slots/" + entry.Name())
		if err != nil {
			panic(err)
		}
		slot := &Slot{}
		if err := json.Unmarshal(filedata, slot); err != nil {
			panic(err)
		}

		ss.store[entry.Name()] = slot
	}
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
	if _, err := ss.putIntoDisk(slot); err != nil {
		return nil, err
	}
	if _, err := ss.putIntoCache(slot); err != nil {
		return nil, err
	}
	return slot, nil
}

func (ss *SlotStore) putIntoCache(slot *Slot) (*Slot, error) {
	if s, ok := ss.store[slot.Id]; ok {
		return s, fmt.Errorf("couldn't create a SLOT. Store already contains ID=%v", s.Id)
	}
	ss.store[slot.Id] = slot
	return slot, nil
}

// putIntoDisk writes slot into disk storage
func (ss *SlotStore) putIntoDisk(slot *Slot) (*Slot, error) {

	json, err := json.Marshal(slot)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(fmt.Sprintf("data/slots/%s.json", slot.Id), json, 0644)
	if err != nil {
		panic(err)
	}
	return slot, nil
}
