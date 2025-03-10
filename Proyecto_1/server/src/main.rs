use std::thread;
use serde::{Deserialize, Serialize};
use std::fs;
use std::io::ErrorKind;
use std::time::Duration;
use serde_json;
use reqwest;

use std::process::Command;

#[derive(Debug, Serialize, Deserialize)]
struct SysInfo {
    mem_total: u64,
    mem_free: u64,
    mem_used: u64,
    cpu_usage: u64,
    containers: Vec<ContainerInfo>,
}

#[derive(Debug, Serialize, Deserialize)]
struct ContainerInfo {
    pid: i32,
    name: String,
    // etc...
}


fn log_conteiner() {
    let nombre_contenedor = "Administrador_logs";

    // Verificar si el contenedor ya existe
    let check = Command::new("docker")
        .args(["ps", "-a", "--format", "{{.Names}}"])
        .output()
        .expect("Error al ejecutar docker ps");

    let contenedor_existente = String::from_utf8_lossy(&check.stdout);

    if contenedor_existente.contains(nombre_contenedor) {
        println!("El contenedor '{}' ya existe. No es necesario crearlo.", nombre_contenedor);
    } else {
        println!("Creando el contenedor '{}'...", nombre_contenedor);

        let status = Command::new("docker")
            .args([
                "run", "-d", "--name", nombre_contenedor,
                "-v", "/var/log:/logs",
                "python:3.9", "tail", "-f", "/logs/syslog"
            ])
            .status()
            .expect("Error al ejecutar docker run");

        if status.success() {
            println!("Contenedor '{}' creado exitosamente.", nombre_contenedor);
        } else {
            println!("Error al crear el contenedor '{}'.", nombre_contenedor);
        }
    }
}

fn set_crontab(client: &Client, action: i32) -> Result<(), Box<dyn Error>> {
    let response = client.post("http://127.0.0.1:8000/setCronJob")
        .json(&action)
        .send()?;

    println!("Respuesta: {}", response.text()?);    
    Ok(())
}


fn main() -> Result<(), Box<dyn std::error::Error>> {
    //log_conteiner();
    let cliente = reqwest::blocking::Client::new();
    set_crontab(cliente,1); // crear el cronjob
    loop {  // Romper bucle con ctrl+c        
        // Leer /proc/sysinfo_<carnet>
        let contenido = match fs::read_to_string("/proc/sysinfo_202112345") {
            Ok(texto) => texto,
            Err(e) if e.kind() == ErrorKind::NotFound => {
                println!("Archivo no encontrado, continuando...");

                continue; 
            },
            Err(e) => {
                println!("Error al leer el archivo: {:?}", e);

                continue; 
            }
        };
        
        let sys_info: SysInfo = match serde_json::from_str(&contenido) {
            Ok(info) => info,
            Err(e) => {
                println!("Error al parsear JSON: {:?}", e);
                continue;
            }
        };

        let sys_info_json = match serde_json::to_string(&sys_info) {
            Ok(json) => json,
            Err(e) => {
                println!("Error al serializar JSON: {:?}", e);

                continue;
            }
        };

        //     - Ver contenedores activos
        //     - Comparar con la regla de "debe haber 1 contenedor de cada tipo"
        //     - Eliminar contenedores sobrantes o viejos
        
        
        let response = match cliente.post("http://127.0.0.1:8000/logs")
            .header("Content-Type", "application/json")
            .body(sys_info_json)
            .send()
        {
            Ok(resp) => resp,
            Err(e) => {
                println!("Error en la petición HTTP: {:?}", e);

                continue;
            }
        };

        println!("Response: {:?}", response);
        println!("Memoria usada: {} MB", sys_info.mem_used);

        thread::sleep(Duration::from_secs(10)); // delay 10 segundos para no saturar
    }
}