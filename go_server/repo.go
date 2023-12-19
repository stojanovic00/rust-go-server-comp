package main

import "sync"

type Repo struct {
	Entries map[int64]MapEntry
}

func NewRepo() *Repo {
	return &Repo{
		Entries: make(map[int64]MapEntry),
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
