use std::sync::{Arc, mpsc, Mutex};
use std::thread;

//Worker thread wrapper
#[allow(dead_code)] //Because of id
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

           //TODO make config param print debug
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


// Job is type that defines closure that can be sent across threads safely
// dyn concrete type is not known during compile time
// FnOnce -> closure can be called only once
// Send -> closure can be sent safely across threads
// 'static - closure is valid during whole program lifetime
type Job = Box<dyn FnOnce() + Send + 'static>;

pub struct ThreadPool{
    workers: Vec<Worker>,
    job_dispatcher: Option<mpsc::Sender<Job>>
}

impl ThreadPool{
    pub fn new(size: usize) -> Self{
        assert!(size > 0);

        let (job_dispatcher, job_listener) = mpsc::channel();

        let job_listener = Arc::new(Mutex::new(job_listener));

        let mut workers = Vec::with_capacity(size);

        for id in 0..size{
            workers.push(Worker::new(id, Arc::clone(&job_listener)));
        }

        ThreadPool{
            workers,
            job_dispatcher: Some(job_dispatcher)
        }
    }

    //Wraps passed closure in Box and dispatches it to workers
   pub fn execute<F>(&self, func: F)
   where
        //Fulfills Job type
        F: FnOnce() + Send + 'static
   {
       let job = Box::new(func);

       self.job_dispatcher.as_ref().unwrap().send(job).unwrap();
   }
}

impl Drop for ThreadPool{
    fn drop(&mut self) {
        drop(self.job_dispatcher.take());

        for worker in &mut self.workers{
            //Waits for currently running jobs in threads to finish
            if let Some(thread) = worker.thread.take(){
                thread.join().unwrap();
            }
        }

    }
}


