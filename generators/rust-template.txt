#![no_main]

use extism_pdk::*;
use serde::{Serialize, Deserialize};


// return a hashmap of array of strings
fn headers_map() ->     std::collections::HashMap<String, Vec<String>> {
    let mut headers = std::collections::HashMap::new();
    headers.insert("Content-Type".to_string(), vec!["application/json".to_string()]);
    headers
}

#[derive(Serialize)]
struct ResponseData {
    pub body: String,
    pub header: std::collections::HashMap<String, Vec<String>>,
    pub code: i32,
}

#[derive(Serialize, Deserialize)]
struct RequestData {
    pub body: String,
    pub header: std::collections::HashMap<String, Vec<String>>,
    pub method: String,
    pub uri: String,
}

#[plugin_fn]
pub fn handle(input: String) -> FnResult<Json<ResponseData>> {

    let req: RequestData = serde_json::from_str(&input).unwrap();

    let message: String = "👋 Hello ".to_string() + &req.body;

    let resp = ResponseData { 
        body: message , 
        code: 200, 
        header: headers_map()
    };
    
    Ok(Json(resp))
}
