struct Worker {
    id: usize,
    thread: Option<thread::JoinHandle<()>>,
}

impl Worker{
    fn new(id: usize, job_listener: Arc<Mutex<mpsc::Receiver<Job>>>) -> Worker{
       //First we spin the thread and then assign its handle to worker

       let thread = thread::spawn(move || loop {
           //All threads will wait for their turn to lock the channel and by that receive some job
          let message = job_listener.lock().unwrap().recv();

           match message {
              Ok(job ) => {
                  println!("Worker {id} got a job; executing.");
                  job();
              },
               //Error will occur when job_dispatcher is closed
               Err(_) => {
                   println!("Worker {id} disconnected; shutting down.");
                   break;
               }
           }
       });

        Worker{
            id,
            thread: Some(thread)
        }
    }
}

