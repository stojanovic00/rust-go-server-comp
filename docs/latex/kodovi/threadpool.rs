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
}


