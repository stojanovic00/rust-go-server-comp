for i := 0; i < poolSize; i++ {
  wg.Add(1)
  go func(threadNum int) {
    defer wg.Done()

    for {
      select {
      case conn, ok := <-connChan:
        {
          if !ok {
            return
          }
          fmt.Printf("Thread %d handles request\n", threadNum)
          handleConnection(conn, repo, mapMux)
        }
      }
    }
  }(i)
}

