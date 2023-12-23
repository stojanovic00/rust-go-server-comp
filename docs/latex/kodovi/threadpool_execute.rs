impl ThreadPool{
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
