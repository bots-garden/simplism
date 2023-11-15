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
struct ReturnValue {
    pub body: String,
    pub header: std::collections::HashMap<String, Vec<String>>,
    pub code: i32,
}

#[derive(Serialize, Deserialize)]
struct Argument {
    pub body: String,
    pub header: std::collections::HashMap<String, Vec<String>>,
    pub method: String,
    pub uri: String,
}


#[plugin_fn]
pub fn hello(_: String) -> FnResult<Json<ReturnValue>> {

    let message: String = "👋 Hello World 🌍".to_string();

    let output = ReturnValue { 
        body: message , 
        code: 200, 
        header: headers_map()
    };
    
    Ok(Json(output))
}
