#!/bin/bash
# Este script crea 10 contenedores aleatorios utilizando la imagen containerstack/alpine-stress
# Seleccion aleatoria de estres: cpu, mem, io y disk.
# cronjob a ejecutar cada 30 segundos
# docker container prune -f <-- Contenedores detenidos

# Número de contenedores a crear

NUM_CONTAINERS=10
Tipo_Estress=("cpu" "mem" "io" "disk")

Tiempo_de_vida=30

for i in $(seq 1 $NUM_CONTAINERS); do

    Nombre_Contenedor="stress_$(date +%s%N)" # Nombre unico con fecha y nanosegundos
    
    Tipo_Estres=$(printf "%s\n" "${Tipo_Estress[@]}" | shuf -n 1) # Aleatorios
    
    STRESS_DURATION=30 #duracion estres
    
    case "$Tipo_Estres" in
        cpu)
            Docker_Cmd="--cpu 1 --timeout ${STRESS_DURATION}s"
            ;;
        mem)
            Docker_Cmd="--vm 1 --vm-bytes 100M --timeout ${STRESS_DURATION}s"
            ;;
        io)
            Docker_Cmd="--io 1 --timeout ${STRESS_DURATION}s"
            ;;
        disk|*)
            Docker_Cmd="--hdd 1 --timeout ${STRESS_DURATION}s"
            ;;
    esac

    echo "Creando contenedor '$Nombre_Contenedor' de tipo '$Tipo_Estres'..."
    docker run -d --name "$Nombre_Contenedor" \
    --memory="128m"  --cpus="0.2" \
    containerstack/alpine-stress stress $(echo $Docker_Cmd)
    
# --memory-swap="128m"

    if [ $? -eq 0 ]; then
        echo "-----------------$i-------------------"
        echo "Contenedor $Nombre_Contenedor creado exitosamente."

        (sleep $Tiempo_de_vida && docker stop "$Nombre_Contenedor" && docker rm "$Nombre_Contenedor") & # Programar eliminacion
        echo "-----------------$i-------------------"
    else
        echo "Error al crear el contenedor $Nombre_Contenedor."
    fi
done

echo "Creación de $NUM_CONTAINERS contenedores completada."
