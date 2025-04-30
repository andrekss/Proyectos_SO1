use actix_web::{get, post, web, App, HttpServer, Responder};
use serde::{Deserialize, Serialize};
use std::sync::Mutex;
use once_cell::sync::Lazy;

#[derive(Deserialize, Serialize, Debug, Clone)]
struct Tweet {
    description: String,
    country: String,
    weather: String,
}

static tweet_recibido: Lazy<Mutex<Option<Tweet>>> = Lazy::new(|| Mutex::new(None)); // global

#[post("/input")]
async fn handle_tweet(tweet: web::Json<Tweet>) -> impl Responder {
    println!("Tweet recibido: {:?}", tweet);

    // Guardar tweet
    let mut guard = tweet_recibido.lock().unwrap();
    *guard = Some(tweet.into_inner());
    
    "Tweet subido"
}

#[get("/get_tweet")]
async fn get_tweet() -> impl Responder {
    let guard = tweet_recibido.lock().unwrap();
    if let Some(tweet) = &*guard {
        actix_web::HttpResponse::Ok()
            .content_type("application/json")
            .json(tweet) // conversor a json
    } else {
        actix_web::HttpResponse::NotFound()
            .content_type("application/json")
            .body("\"No hay tweet guardado\"")
    }
}


#[get("/")]
async fn funcionamiento() -> impl Responder {
    "Servicio funcionando . . ."
}

#[get("/health")]
async fn health_check() -> impl Responder {
    actix_web::HttpResponse::Ok().body("OK")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Servidor funcionando...");
    HttpServer::new(|| {
        App::new()
            .service(handle_tweet)
            .service(get_tweet)
            .service(funcionamiento)
            .service(health_check)
    })
    .bind(("0.0.0.0", 8082))?
    .run()
    .await
}
