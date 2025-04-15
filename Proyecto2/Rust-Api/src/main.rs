use actix_web::{post, web, App, HttpServer, Responder};
use serde::Deserialize;

#[derive(Deserialize, Debug)]
struct Tweet {
    description: String,
    country: String,
    weather: String,
}

#[post("/input")]
async fn handle_tweet(tweet: web::Json<Tweet>) -> impl Responder {
    println!("Tweet recibido: {:?}", tweet);
    "Tweet subido"
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| App::new().service(handle_tweet))
        .bind(("0.0.0.0", 8080))?
        .run()
        .await
}
