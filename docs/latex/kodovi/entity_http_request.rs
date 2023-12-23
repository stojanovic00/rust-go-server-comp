pub struct Entity{
    pub id: i64,
    pub description: String,
    pub value: f32
}

pub struct HttpRequest{
    pub method: HttpMethod,
    pub path: String,
    pub headers: Vec<String>,
    pub body: Option<Entity>
}


