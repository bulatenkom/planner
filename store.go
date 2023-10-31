package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Ider interface {
	GetId() string
}

// L1 Cache
type Store[T Ider] struct {
	store            map[string]*T
	entityName       string // entity name of T
	entityNamePlural string
	entityLocation   string
}

func NewStore[T Ider]() Store[T] {
	ss := Store[T]{}
	ss.store = make(map[string]*T, 1024)
	ss.entityName = strings.ToLower(reflect.TypeOf(new(T)).Elem().Name())
	ss.entityNamePlural = ss.entityName + "s"
	ss.entityLocation = filepath.Join(AppFlags.DataRoot, ss.entityNamePlural)
	// ------------------------------------------------------------
	// Init directories
	// ------------------------------------------------------------

	if err := os.MkdirAll(ss.entityLocation, 0774); err != nil {
		panic(err)
	}
	// ------------------------------------------------------------
	// load file into store
	// ------------------------------------------------------------
	dirEntries, err := os.ReadDir(ss.entityLocation)
	if err != nil {
		panic(err)
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		filedata, err := os.ReadFile(filepath.Join(ss.entityLocation, entry.Name()))
		if err != nil {
			panic(err)
		}

		entity := new(T)
		if err := json.Unmarshal(filedata, entity); err != nil {
			panic(err)
		}

		ss.store[entry.Name()] = entity
	}
	return ss
}

func (ss *Store[T]) FindAll() map[string]*T {
	return ss.store
}

func (ss *Store[T]) FindById(id string) (*T, error) {
	if entity, ok := ss.FindAll()[id]; ok {
		return entity, nil
	} else {
		return nil, fmt.Errorf("could not find %s with id=%q", ss.entityName, id)
	}
}

func (ss *Store[T]) Create(entity *T) (*T, error) {
	if _, err := ss.putIntoDisk(entity); err != nil {
		return nil, err
	}
	if _, err := ss.putIntoCache(entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (ss *Store[T]) putIntoCache(entity *T) (*T, error) {
	if s, ok := ss.store[(*entity).GetId()]; ok {
		return s, fmt.Errorf("couldn't create a %s. Store already contains ID=%v", ss.entityName, (*s).GetId())
	}
	ss.store[(*entity).GetId()] = entity
	return entity, nil
}

// putIntoDisk writes entity into disk storage
func (ss *Store[T]) putIntoDisk(entity *T) (*T, error) {
	json, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(filepath.Join(ss.entityLocation, fmt.Sprintf("%s.json", (*entity).GetId())), json, 0644)
	if err != nil {
		panic(err)
	}
	return entity, nil
}
