use std::io::{BufRead, BufReader, Read};
use std::net::TcpStream;
use crate::model::Entity;


#[derive(Debug)]
pub enum HttpMethod{
    Get,
    Put,
    NotImplemented
}
#[derive(Debug)]
pub struct HttpRequest{
    pub method: HttpMethod,
    pub path: String,
    pub headers: Vec<String>,
    pub body: Option<Entity>
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

pub fn parse_request(stream: &mut TcpStream) -> Result<HttpRequest, String>{
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
        if line.to_lowercase().contains("content-length"){
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
