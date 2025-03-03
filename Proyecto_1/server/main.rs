use std::thread;
use std::time::Duration;
use serde::{Deserialize, Serialize};
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

fn main() -> Result<(), Box<dyn std::error::Error>> {
    log_conteiner();

    loop {  // Romper bucle con ctrl+c

        // Leer /proc/sysinfo_<carnet>
        let contenido = std::fs::read_to_string("/proc/sysinfo_202112345")?;
        let sys_info: SysInfo = serde_json::from_str(&contenido)?;

        //     - Ver contenedores activos
        //     - Comparar con la regla de "debe haber 1 contenedor de cada tipo"
        //     - Eliminar contenedores sobrantes o viejos
        
        let cliente = reqwest::blocking::Client::new();
        client.post("http://127.0.0.1:8000/logs")
           .json(&sys_info)  // enviamos el json
           .send()?;


        println!("Memoria usada: {} MB", sys_info.mem_used);

        thread::sleep(Duration::from_secs(10)); // delay 10 segundos para no saturar
    }

    // al salir capturar ctrl+c para cleanup
    // petición generar gráficas
    // Eliminar cronjob
    // mas...
}
