use std::thread;
use std::fs;
use std::io::ErrorKind;
use std::time::Duration;
use std::error::Error;
use serde_json::{self, Value, json}; // Importamos `json!`
use chrono::Utc; // Importamos `Utc`
use reqwest::blocking::Client;
use std::sync::{Arc, atomic::{AtomicBool, Ordering}};
use std::process::{Command, Output};
use std::collections::HashMap;
use ctrlc;



fn log_conteiner(action: i32) -> Result<(), Box<dyn Error>> {
    match action {
        1 => {
            let status = Command::new("sh")
                .arg("/home/andres/Escritorio/Sistemas_Operativos_1/Proyectos_SO1/Proyecto_1/Logs/deploy.sh") // Ajusta la ruta a tu script
                .status()?;
            if status.success() {
                println!("Script ejecutado correctamente");
            } else {
                println!("Hubo un error al ejecutar el script");
            }
        },
        0 => {
            let output = Command::new("docker")
            .args(["stop", "Servicio"])
            .output()?; // Capturamos stdout y stderr
        
        if output.status.success() {
            println!("Se paró el servicio.");
        } else {
            println!("Error al ejecutar el comando: {}", 
                String::from_utf8_lossy(&output.stderr)); // Mostramos el error de Docker
        }
        
        },
        _ => {
            println!("Use 1 o 0.");
        }
    }
    Ok(())
}

fn Ejecutar_Modulo_Kernel() -> Result<(), Box<dyn Error>> {
    let status = Command::new("sh")
        .arg("/home/andres/Escritorio/Sistemas_Operativos_1/Proyectos_SO1/Proyecto_1/module/Exe.sh") // Ajusta la ruta a tu script
        .status()?;
    if status.success() {
        println!("Modulo kernel ejecutado correctamente");
    } else {
        println!("Hubo un error al ejecutar el Modulo kernel");
    }
    Ok(())
}

fn Limpiar_Modulo_Kernel() -> Result<(), Box<dyn Error>> {
    let status = Command::new("sh")
        .arg("/home/andres/Escritorio/Sistemas_Operativos_1/Proyectos_SO1/Proyecto_1/module/Eliminar_Modulo.sh") // Ajusta la ruta a tu script
        .status()?;
    if status.success() {
        println!("Modulo kernel Limpiado Correctamente");
    } else {
        println!("Hubo un error al ejecutar el la limpieza del modulo kernel");
    }
    Ok(())
}

fn get_sysinfo_json() -> Result<String, Box<dyn Error>> {
    let contenido = match fs::read_to_string("/proc/sysinfo_202113580") {
        Ok(texto) => texto,
        Err(e) if e.kind() == ErrorKind::NotFound => return Err("Archivo no encontrado".into()),
        Err(e) => return Err(e.into()),
    };

    let json_value: Value = serde_json::from_str(&contenido)?;

    let sys_info_json = serde_json::to_string_pretty(&json_value)?;
    Ok(sys_info_json)
}


fn set_crontab(action: i32) -> Result<(), Box<dyn Error>> {
    match action {
        1 => {
            let lines = "\
* * * * * /home/andres/Escritorio/Sistemas_Operativos_1/Proyectos_SO1/Proyecto_1/Scripts/Contenedores.sh
* * * * * sleep 30 && /home/andres/Escritorio/Sistemas_Operativos_1/Proyectos_SO1/Proyecto_1/Scripts/Contenedores.sh";

            let salida = Command::new("crontab").arg("-l").output();
            let mut escribiendo_crontab = match salida {
                Ok(o) if o.status.success() => String::from_utf8_lossy(&o.stdout).to_string(),
                _ => String::new(),
            };

            if !escribiendo_crontab.ends_with('\n') {
                escribiendo_crontab.push('\n');
            }

            escribiendo_crontab.push_str(lines);

            if !escribiendo_crontab.ends_with('\n') {
                escribiendo_crontab.push('\n');
            }

            let cron_file = "/tmp/mycron";
            fs::write(cron_file, &escribiendo_crontab)?;

            let status = Command::new("crontab").arg(cron_file).status()?;
            if !status.success() {
                return Err("Error instalando crontab".into());
            }
            println!("Crontab actualizado con el script cada 30 segundos");
        },
        0 => {
            let check_crontab = Command::new("crontab").arg("-l").output()?;
            if !check_crontab.status.success() {
                println!("No hay crontab activo para borrar.");
                return Ok(());
            }

            let status = Command::new("crontab").arg("-r").status()?;
            if !status.success() {
                println!("Error borrando crontab");
            } else {
                println!("Crontab borrado exitosamente.");
            }
        },
        _ => {
            return Err("1 para crear, 0 para borrar".into());
        }
    }
    Ok(())
}


fn limpiar_contenedores() -> Result<(), Box<dyn Error>> {
    let output: Output = Command::new("docker")
        .args(["ps", "--format", "{{.ID}} {{.Names}} {{.Command}}"])
        .output()?;
    
    if !output.status.success() {
        println!("Error al obtener la lista de contenedores");
        return Ok(());
    }

    let salida = String::from_utf8_lossy(&output.stdout); // Obtenemos las líneas

    let mut contenedores_por_tipo: HashMap<String, Vec<String>> = HashMap::new();

    // Recorre cada línea de "<ID> <NOMBRE> <COMMAND>" para parts
    for line in salida.lines() {
        let parts: Vec<&str> = line.split_whitespace().collect();
        if parts.len() < 3 {
            continue;
        }

        let container_id = parts[0].to_string(); // id
        let command_str = &line[line.find(parts[2]).unwrap_or(0)..]; // command

        let categoria = if command_str.contains("--cpu") {
            "cpu"
        } else if command_str.contains("--vm") {
            "mem"
        } else if command_str.contains("--io") {
            "io"
        } else if command_str.contains("--hdd") {
            "disk"
        } else { // no coincide
            continue;
        };

        contenedores_por_tipo
            .entry(categoria.to_string())
            .or_insert_with(Vec::new)
            .push(container_id);
    }

    for (tipo, ids) in contenedores_por_tipo.iter() {
        if ids.len() > 1 {
            // Mantenemos el primero, eliminamos el resto
            for id in &ids[1..] {
                println!("Contenedor eliminado extra de tipo {}: {}", tipo, id);

                // Detenemos y removemos
                let _ = Command::new("docker")
                    .args(["stop", id])
                    .output();

                let _ = Command::new("docker")
                    .args(["rm", id])
                    .output();
            }
        }
    }

    Ok(())
}


fn Deserializar_Json_Y_Formatear(json_str: &str) -> Result<Value, Box<dyn std::error::Error>> {
    let parsed_json: Value = serde_json::from_str(json_str)?;

    // Extraer los procesos
    let procesos = parsed_json["Processes"]
        .as_array()
        .ok_or("No se encontraron procesos en el JSON")?;

    let mut registros = Vec::new();

    

    // Recorrer cada proceso y agregarlo como un nuevo registro
    for proceso in procesos {
        let timestamp = Utc::now().to_rfc3339(); // generamos timestamp
        let registro = json!({
            "time": timestamp,
            "CPU": proceso["Porcentaje de uso CPU"].as_f64().unwrap_or(0.00),
            "Mem": proceso["Porcentaje de uso Memoria"].as_f64().unwrap_or(0.00),
            "IO": proceso["Operaciones I/O"].as_i64().unwrap_or(0),
            "Disc": proceso["Uso de disco"].as_i64().unwrap_or(0)
        });
        registros.push(registro);
        thread::sleep(Duration::from_millis(10));
    }

    Ok(Value::Array(registros))
}


fn main() -> Result<(), Box<dyn std::error::Error>> {

    let cliente = reqwest::blocking::Client::new();
    
    let corriendo = Arc::new(AtomicBool::new(true));{ // Controlar ciclo
        let corriendo = corriendo.clone();
        let cliente_clone = cliente.clone(); // Clonamos el cliente para usarlo en el handler
        ctrlc::set_handler(move || {
            println!("Parando ciclo . . .");
            let _ = set_crontab(0); // borramos crontab
            corriendo.store(false, Ordering::SeqCst);
        })?;
    }
    set_crontab(1)?; // crear el cronjob para generar cronjob
    log_conteiner(1)?; // Creamos el contenedor con el servicio python
    Ejecutar_Modulo_Kernel()?; 
    while corriendo.load(Ordering::SeqCst) {  // Romper bucle con ctrl+c     
        
        limpiar_contenedores();
        match get_sysinfo_json() {
            Ok(json_str) => {
                match Deserializar_Json_Y_Formatear(&json_str) {
                    Ok(json_Formateado) => println!("{}", serde_json::to_string_pretty(&json_Formateado).unwrap()),
                    Err(e) => println!("Error transformando JSON: {}", e),
                }
            },
            Err(e) => println!("Error obteniendo JSON: {}", e),
        }
        
// post mandar json actualizado
        //     - Ver contenedores activos
        //     - Comparar con la regla de "debe haber 1 contenedor de cada tipo"
        //     - Eliminar contenedores sobrantes o viejos
      
        /*let response = match cliente.post("http://127.0.0.1:8000/logs")
            .header("Content-Type", "application/json")
            .body(sys_info_json)
            .send()
        {
            Ok(resp) => resp,
            Err(e) => {
                println!("Error en la petición HTTP: {:?}", e);

                continue;
            }
        };*/

        //println!("Response: {:?}", response);
        //println!("Memoria usada: {} MB", sys_info.mem_used);

        thread::sleep(Duration::from_secs(10)); // delay 10 segundos para no saturar
    }
    Limpiar_Modulo_Kernel()?;
    set_crontab(0)?; // borramos crontab al salir
    log_conteiner(0)?; // paramos el contenedor
    Ok(())
}
