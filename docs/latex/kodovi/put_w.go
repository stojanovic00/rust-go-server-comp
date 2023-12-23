{
  ...
  //Entry doesn't exist -> lock whole map
  mapMux.Lock()
  repo.Entries[request.Body.Id] = *NewMapEntry(*request.Body)
  mapMux.Unlock()
  ...
}

