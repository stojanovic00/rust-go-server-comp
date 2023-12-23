impl Repo{
    pub fn get_by_id(&self, id: i64) ->  Option<Entity>{
        match self.entities.get(&id){
          Some(entity_lock) => {
              let entity = entity_lock.read().ok()?;

              Some((*entity).clone())
           },
          None => None
        }
    }
}

{
...
if let Ok(ro_repo) = repo.read(){
     match ro_repo.get_by_id(id){
        Some(entity) => {
            let json_entity = serde_json::to_string(&entity).unwrap();
            let ok_resp = "HTTP/1.1 200 OK \r\n";
            let content_len = json_entity.len();

            let response = format!(
                "{ok_resp}Content-Length: {content_len}\r\nContent-Type: application/json\r\nConnection: close\r\n\r\n{json_entity}"
            );

            stream.write_all(response.as_bytes()).unwrap();
        }
        None =>{
            let not_found_resp = "HTTP/1.1 404 NotFound \r\n\r\n";
            stream.write_all(not_found_resp.as_bytes()).unwrap();
            return;
        }
    };
}
...
}

