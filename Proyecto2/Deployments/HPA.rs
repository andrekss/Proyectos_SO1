// rust_api/src/main.rs
use actix_web::{post, web, App, HttpServer, Responder};
use serde::Deserialize;

#[derive(Deserialize)]
struct Tweet {
    Description: String,
    Country: String,
    Weather: String,
}

#[post("/input")]
async fn handle_tweet(tweet: web::Json<Tweet>) -> impl Responder {
    println!("Received tweet: {:?}", tweet);
    // Aquí llamarías a la API REST o gRPC de Go
    "Tweet forwarded!"
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| App::new().service(handle_tweet))
        .bind(("0.0.0.0", 8080))?
        .run()
        .await
}
