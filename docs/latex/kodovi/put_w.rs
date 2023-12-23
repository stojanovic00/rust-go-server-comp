{
...
// Doesnt exist: lock whole map and add
let mut w_repo = repo.write().unwrap();
let new_id = body.id;
let new_ent =  Arc::new(RwLock::new(body));
w_repo.entities.insert(new_id, new_ent);
...
}
