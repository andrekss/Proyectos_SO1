#!/bin/bash

HARBOR_IP="34.55.73.46"

# Lista de imágenes que deseas subir
IMAGES=("proyecto2-api-rust" "proyecto2-rabbit-writer" "proyecto2-kafka-writer" "proyecto2-api-go" \
        "proyecto2-kafka-consumer" "proyecto2-rabbit-consumer" "bitnami/zookeeper" "redis" \
        "valkey/valkey" "grafana/grafana" "bitnami/kafka" "rabbitmq")

# Iniciar sesión en Harbor
#docker login $HARBOR_IP

# Subir imágenes
for IMAGE in "${IMAGES[@]}"; do
    echo "Subiendo la imagen: $IMAGE"
    
    docker tag $IMAGE $HARBOR_IP/library/$IMAGE

    docker push $HARBOR_IP/library/$IMAGE
done

echo "¡Todas las imágenes han sido subidas a Harbor!"
