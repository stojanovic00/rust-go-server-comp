use std::{net::{TcpListener, TcpStream}, collections::HashMap};
use  std::io::{prelude::*};
use std::sync::{Arc, Mutex, RwLock};

use ctrlc;
use dotenv::dotenv;


mod thread_pool;
mod http_utils;
mod model;

use thread_pool::ThreadPool;
use  http_utils::{HttpRequest, HttpMethod};
use  model::Entity;


pub struct Repo{
    entities: HashMap<i64,Arc<RwLock<Entity>>>
}

impl Repo{
    pub fn new() -> Self{
        Self{
            entities:HashMap::new()
        }
    }

    pub fn get_by_id(&self, id: i64) ->  Option<Entity>{
        //Now accessing entry level RwLock
        match self.entities.get(&id){
          Some(entity_lock) => {
              // Read from the RwLock
              let entity = entity_lock.read().ok()?;

              Some((*entity).clone())
           },
          None => None
        }
    }
}


fn main() {
    // Load environment variables from the .env file
    dotenv().ok();

    let tcp_address = std::env::var("TCP_ADDRESS").expect("TCP_ADDRESS not set in .env");

    let thread_pool_size = std::env::var("THREAD_POOL_SIZE").expect("THREAD_POOL_SIZE not set in .env");
    let thread_pool_size: usize = thread_pool_size.parse().expect("Invalid THREAD_POOL_SIZE value");

    println!("Server started on address: {}", tcp_address);
    println!("Thread pool size: {}", thread_pool_size);


    //Init resources
    let listener = TcpListener::bind(tcp_address).unwrap();

    //Allowing multiple reads and single write
    let repo = Arc::new(RwLock::new(Repo::new()));

    let thread_pool = ThreadPool::new(5);



    // Use a Mutex to track whether to continue the loop
    let terminated = Arc::new(Mutex::new(false));

    // Listening for CTRL + C signal to perform graceful shutdown
    let terminated_clone = Arc::clone(&terminated);
    ctrlc::set_handler(move || {
        println!("\nServer stop signal received");
        println!("Shutting down server...");

        // Set the flag to stop the loop
        let mut termintd = terminated_clone.lock().unwrap();
        *termintd = true;

        //Sends tcp request to make sure there will be one more iteration inside loop that listens for tcp connections
         TcpStream::connect("localhost:7878").unwrap();
    })
    .expect("Error setting Ctrl+C handler");


    // Listening for TCP
    let termntd_clone = Arc::clone(&terminated);
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
}

fn handle_connection(mut stream: TcpStream, repo: Arc<RwLock<Repo>>){
    let request =http_utils::parse_request(&mut stream);

    match request{
        Ok(request) => {
            match &request.method {
                HttpMethod::Get => handle_get(request, &mut stream, repo),
                HttpMethod::Put => handle_put(request, &mut stream,  repo),
                HttpMethod::NotImplemented => {
                    let response = "HTTP/1.1 501 NotImplemented \r\n\r\n";
                    stream.write_all(response.as_bytes()).unwrap();
                    return;
                }
            }
        },
        Err(err) => {
            stream.write_all(err.as_bytes()).unwrap();
        }

    }
}


fn handle_get(request: HttpRequest ,stream: &mut TcpStream, repo: Arc<RwLock<Repo>>){
    let id:i64;

    //Extract id
    match request.path.strip_prefix("/") {
       Some(number_str) => {
            if let Ok(number) = number_str.parse::<i64>(){
                id = number;
            } else{
                let response = "HTTP/1.1 400 BadRequest \r\n\r\n";
                stream.write_all(response.as_bytes()).unwrap();
                return;
            }

       }, 
       None => {
            let response = "HTTP/1.1 400 BadRequest \r\n\r\n";
            stream.write_all(response.as_bytes()).unwrap();
            return;
       }
    } 


    if let Ok(ro_repo) = repo.read(){
         match ro_repo.get_by_id(id){
            Some(entity) => {
                let json_entity = serde_json::to_string(&entity).unwrap();
                let ok_resp = "HTTP/1.1 200 OK \r\n";
                let content_len = json_entity.len();

                let response = format!(
                    "{ok_resp}Content-Length: {content_len}\r\nContent-Type: application/json\r\n\r\n{json_entity}"
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
}



fn handle_put(request: HttpRequest, stream: &mut TcpStream, repo: Arc<RwLock<Repo>>){
    if request.path != "/"{
        let response = "HTTP/1.1 400 BadRequest \r\n\r\n";
        stream.write_all(response.as_bytes()).unwrap();
        return;
    }
    if let Some(body)  = request.body{
        let exists: bool;

        //In new scope to ensure unlock
        {
            let r_repo = repo.read().unwrap();
            exists = r_repo.entities.contains_key(&body.id);
        }

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
            // Doesnt exist: lock whole map and add
            let mut w_repo = repo.write().unwrap();
            let new_id = body.id;
            let new_ent =  Arc::new(RwLock::new(body));
            w_repo.entities.insert(new_id, new_ent);
        }



        let response = "HTTP/1.1 200 OK \r\n\r\n";
        stream.write_all(response.as_bytes()).unwrap();
        return;
    } else {
        let response = "HTTP/1.1 400 BadRequest \r\n\r\n";
        stream.write_all(response.as_bytes()).unwrap();
        return;
    }
}