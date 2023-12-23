{
...
mapMux.RLock()
entry, ok := repo.Entries[id]


entry.Mux.RLock()
jsonBytes, err := json.Marshal(entry.Entity)
entry.Mux.RUnlock()

mapMux.RUnlock()
...
}
