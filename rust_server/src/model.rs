use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Entity{
    pub id: i64,
    pub description: String,
    pub value: f32
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
