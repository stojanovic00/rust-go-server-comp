{
...
if exists {
    //just mutate entry without locking whole map
    if let Some(ro_repo) = repo.read().ok() {
        let entry = ro_repo.entities.get(&body.id).unwrap();

        //Locking and changing entity inside repo
        if let Some(mut rw_entity) = entry.write().ok() {
            *rw_entity = body;
        }else{
            let response = "HTTP/1.1 500 InternalServerError \r\n\r\n";
            stream.write_all(response.as_bytes()).unwrap();
            return;
        }
    } else {
        let response = "HTTP/1.1 500 InternalServerError \r\n\r\n";
        stream.write_all(response.as_bytes()).unwrap();
        return;
    }
} else {
...
}
