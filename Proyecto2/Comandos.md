gcloud auth login
gcloud container clusters list

// Conectar al clúster
gcloud container clusters get-credentials so-cluster-1 \
  --zone us-central1-c \
  --project reliable-byte-455302-k8


sudo du -h / --max-depth=2 | sort -hr | head -20 // limpieza
docker-compose up --build // subir y construir

// crear topic
docker exec -it kafka /opt/bitnami/kafka/bin/kafka-topics.sh \
  --create \
  --topic tweets \
  --bootstrap-server localhost:9092 \
  --partitions 1 \
  --replication-factor 1


// verificar existencia
docker exec -it kafka /opt/bitnami/kafka/bin/kafka-topics.sh \
  --list \
  --bootstrap-server localhost:9092


docker rmi -f $(docker images -aq) // borra todo cache de imagenes

// consumers
go mod init multi_subscriber
go get github.com/segmentio/kafka-go
go get github.com/streadway/amqp
go mod tidy

// eliminar todo docker
docker system prune -a --volumes


// base de datos
docker run -d --name redis -p 6379:6379 redis:latest

// libreria go con redis

go get github.com/go-redis/redis/v8

go get github.com/go-redis/redis/v8

