for stream in listener.incoming() {
    //Check for termination
    let terminate = termntd_clone.lock().unwrap();
    if *terminate {
        break;
    }

    let stream = stream.unwrap();
    let repo_instance = Arc::clone(&repo);

    thread_pool.execute(move || {
        handle_connection(stream, repo_instance);
    })

}

