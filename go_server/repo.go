package main

import "sync"

type Repo struct {
	Entries sync.Map
}

func NewRepo() *Repo {
	return &Repo{
		Entries: sync.Map{},
	}
}

type MapEntry struct {
	Mux    *sync.RWMutex
	Entity Entity
}

func NewMapEntry(entity Entity) *MapEntry {
	return &MapEntry{
		Entity: entity,
		Mux:    &sync.RWMutex{},
	}
}
func NewMapEntryWMux(entity Entity, mux *sync.RWMutex) *MapEntry {
	return &MapEntry{
		Entity: entity,
		Mux:    mux,
	}
}
