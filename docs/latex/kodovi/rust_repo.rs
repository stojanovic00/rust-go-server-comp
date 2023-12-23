pub struct Repo{
    entities: HashMap<i64,Arc<RwLock<Entity>>>
}

let repo = Arc::new(RwLock::new(Repo::new()));


