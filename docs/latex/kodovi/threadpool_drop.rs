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
