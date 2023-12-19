package main

type Repo struct {
	Entities map[int64]Entity
}

func NewRepo() *Repo {
	return &Repo{
		Entities: make(map[int64]Entity),
	}
}
