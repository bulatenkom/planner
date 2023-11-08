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

		ss.store[(*entity).GetId()] = entity
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
	if ss.exists(entity) {
		return entity, fmt.Errorf("[Collision occurred] couldn't create a %s. Store already contains ID=%v", ss.entityName, (*entity).GetId())
	}
	if _, err := ss.putIntoDisk(entity); err != nil {
		return nil, err
	}
	if _, err := ss.putIntoCache(entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (ss *Store[T]) Update(entity *T) (*T, error) {
	if !ss.exists(entity) {
		return entity, fmt.Errorf("couldn't update a %[1]s. Store does not contain %[1]s with ID=%v", ss.entityName, (*entity).GetId())
	}
	if _, err := ss.putIntoDisk(entity); err != nil {
		return nil, err
	}
	if _, err := ss.putIntoCache(entity); err != nil {
		// we don't revert disk write to not lose user data, but it can affect next App startup which will require issue investigation
		return nil, err
	}
	return entity, nil
}

// putIntoCache writes entity into L1 cache (if exists then override)
func (ss *Store[T]) putIntoCache(entity *T) (*T, error) {

	ss.store[(*entity).GetId()] = entity
	return entity, nil
}

func (ss *Store[T]) exists(entity *T) bool {
	_, ok := ss.store[(*entity).GetId()]
	return ok
}

func (ss *Store[T]) existsById(id string) bool {
	_, ok := ss.store[id]
	return ok
}

// putIntoDisk writes entity into disk storage (if exists then override)
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

func (ss *Store[T]) deleteById(id string) error {
	if !ss.existsById(id) {
		return fmt.Errorf("couldn't delete a %[1]s. Store does not contain %[1]s with ID=%v", ss.entityName, id)
	}

	// delete from L1 cache
	entityRef := ss.store[id]
	delete(ss.store, id)
	// delete from disk
	err := os.Remove(filepath.Join(ss.entityLocation, fmt.Sprintf("%s.json", id)))
	if err != nil {
		ss.store[id] = entityRef
		return err
	}
	return nil
}
