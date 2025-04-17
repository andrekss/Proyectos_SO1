use actix_web::{get,post, web, App, HttpServer, Responder};
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

#[get("/")]
async fn funcionamiento() -> impl Responder {
    println!("Servicio funcionando . . .");
    "Servicio funcionando . . ."
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("funcionando . . .");
    HttpServer::new(|| {
        App::new()
        .service(handle_tweet)
        .service(funcionamiento)
        })
        .bind(("0.0.0.0", 8082))?
        .run()
        .await
}
