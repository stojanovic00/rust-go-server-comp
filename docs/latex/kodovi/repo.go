type Repo struct {
	Entries map[int64]MapEntry
}


type MapEntry struct {
	Mux    *sync.RWMutex
	Entity Entity
}

func main(){
...
	repo := NewRepo()
	mapMux := &sync.RWMutex{}
...
}

