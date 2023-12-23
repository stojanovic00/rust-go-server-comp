go func() {
  for {
    conn, err := listener.Accept()
    if err != nil {
      fmt.Println("Error accepting:", err)
      continue
    }

    //Check for termination
    shutdownServerFlag.Lock()
    if shutdownServerFlag.close {
      close(connChan)
      shutdownServerFlag.Unlock()
      return
    } else {
      shutdownServerFlag.Unlock()
    }

    //Dispatch request to thread pool
    connChan <- conn
  }
}()

