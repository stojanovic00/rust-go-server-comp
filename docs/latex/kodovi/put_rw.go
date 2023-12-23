{
...
	if exists {
		mapMux.RLock()
		entry, _ := repo.Entries[request.Body.Id]
		entry.Mux.Lock()
		repo.Entries[request.Body.Id] = *NewMapEntryWMux(*request.Body, entry.Mux)
		entry.Mux.Unlock()
		mapMux.RUnlock()
	} 
...
}
