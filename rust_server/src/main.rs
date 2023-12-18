use std::{net::{TcpListener, TcpStream}, collections::HashMap};
use  std::io::{prelude::*, BufReader};
use std::sync::{Arc, Mutex, RwLock};

use serde::{Serialize, Deserialize};
use ctrlc;


mod thread_pool;
use thread_pool::ThreadPool;

#[derive(Debug)]
enum HttpMethod{
    Get,
    Put,
    NotImplemented
}


#[derive(Debug)]
struct HttpRequest{
        method: HttpMethod,
        path: String,
        headers: Vec<String>,
        body: Option<Entity>
}
impl HttpRequest {
    fn new() -> Self {
        Self{
            method: HttpMethod::NotImplemented,
            path: String::from(""),
            headers: vec![],
            body:None
        }
    }
}


#[derive(Debug, Serialize, Deserialize)]
pub struct Entity{
    id: i64,
    description: String,
    value: f32
}

impl Entity{
    pub fn new(new_id: i64, new_desc: String, new_val: f32) -> Self{
        Self{
            id: new_id,
            description: new_desc,
            value: new_val
        }
    }
}


pub struct Repo{
    entities: HashMap<i64,Entity>
}

impl Repo{
    pub fn new() -> Self{
        Self{
            entities:HashMap::new()
        }
    }

    pub fn get_by_id(&self, id: i64) ->  Option<&Entity>{
        self.entities.get(&id)
    }

    pub fn insert(&mut self, entity: Entity){
        self.entities.insert(entity.id, entity);
    }
}


fn main() {
    //Init resources
    //TODO to config file
    let listener = TcpListener::bind("localhost:7878").unwrap();

    //Allowing multiple reads and single write
    let repo = Arc::new(RwLock::new(Repo::new()));

    //TODO read from env or config file
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
    let request = parse_request(&mut stream);

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

fn parse_request(stream: &mut TcpStream) -> Result<HttpRequest, String>{
    let mut buff_reader = BufReader::new(stream);

    let mut request = HttpRequest::new();

    //Parsing header
    loop {
        let  line = &mut String::from("");
        if let Ok(_) = buff_reader.read_line(line){
            if line == "\r\n"{
                break;
            }
          request.headers.push(line.clone().trim().to_string());
        }
        else {
             let response = "HTTP/1.1 400 BadRequest \r\n\r\n".to_string();
             return Err(response);
        }
    }


    // Parsing request line
    if let Some(request_line) = &request.headers.first(){
        let parts: Vec<&str> = request_line.split(' ').collect();

        request.method =  match parts[0].trim()  {
           "GET" => HttpMethod::Get,
           "PUT" => HttpMethod::Put,
            _ => HttpMethod::NotImplemented
        };

        request.path = parts[1].trim().to_string();
    }
    else{
        let response = "HTTP/1.1 400 BadRequest \r\n\r\n".to_string();
        return Err(response);
    }
    request.headers.remove(0);




    //Parsing body
    //Must be done using content length because otherwise buffer cant find end of line for last line of body because
    //EOF is appended when connection is closed
    let mut content_length = 0;
    for line in &request.headers{
        if line.contains("Content-Length"){
            let parts: Vec<&str> = line.split(":").collect();
            let len_str = parts[1].trim();
            content_length = len_str.parse::<i32>().unwrap();
        }
    }

    if content_length == 0 {
        return Ok(request);
    }

    let mut buffer =  vec![0;content_length as usize];
    let body_content;
    if let Ok(_) = buff_reader.read_exact(&mut buffer){
        body_content =String::from_utf8_lossy(&buffer).to_string();
    }
    else {
        let response = "HTTP/1.1 400 BadRequest \r\n\r\n".to_string();
        return Err(response);
    }

    //Parsing json body to domain model
     match  serde_json::from_str::<Entity>(&body_content){
        Ok(entity) =>{
            request.body = Some(entity);
        },
        Err(err) => {
            let response = format!("HTTP/1.1 400 BadRequest \r\n\r\nJson error: {err}");
            return Err(response);
        }
    };

    Ok(request)
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


    //Read only operation
    //TODO Check lifetime scope for if let
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
    if let Ok(mut w_repo) = repo.write(){
        if let Some(body)  = request.body{
            w_repo.insert(body);
            let response = "HTTP/1.1 200 OK \r\n\r\n";
            stream.write_all(response.as_bytes()).unwrap();
            return;
        }
    }

    let response = "HTTP/1.1 400 BadRequest \r\n\r\n";
    stream.write_all(response.as_bytes()).unwrap();
    return;
}