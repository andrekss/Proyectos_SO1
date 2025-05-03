## Documentación


## Deployments

```
apiVersion: v1
kind: List
items:
#####################################################################
# ZOOKEEPER
#####################################################################
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: zookeeper
    namespace: proyecto2
  spec:
    replicas: 1
    selector: { matchLabels: { app: zookeeper } }
    template:
      metadata: { labels: { app: zookeeper } }
      spec:
        imagePullSecrets:
        - name: harbor-regcred 
        containers:
        - name: zookeeper
          image: 34.55.73.46.nip.io/library/bitnami/zookeeper:latest
          ports: [{containerPort: 2181}]
          env:
            - name: ALLOW_ANONYMOUS_LOGIN
              value: "yes"
- apiVersion: v1
  kind: Service
  metadata:
    name: zookeeper
    namespace: proyecto2
  spec:
    selector: { app: zookeeper }
    ports: [{port: 2181, targetPort: 2181}]
#####################################################################
# KAFKA
#####################################################################
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: kafka
    namespace: proyecto2
  spec:
    replicas: 1
    selector: { matchLabels: { app: kafka } }
    template:
      metadata: { labels: { app: kafka } }
      spec:
        imagePullSecrets:
        - name: harbor-regcred 
        containers:
        - name: kafka
          image: 34.55.73.46.nip.io/library/bitnami/kafka:latest
          ports: [{containerPort: 9092}]
          env:
            - name: KAFKA_CFG_ADVERTISED_LISTENERS
              value: PLAINTEXT://kafka:9092
            - name: KAFKA_CFG_LISTENERS
              value: PLAINTEXT://0.0.0.0:9092
            - name: KAFKA_ZOOKEEPER_CONNECT
              value: zookeeper:2181
            - name: ALLOW_PLAINTEXT_LISTENER
              value: "yes"
            - name: KAFKA_ENABLE_KRAFT
              value: "no"
- apiVersion: v1
  kind: Service
  metadata:
    name: kafka
    namespace: proyecto2
  spec:
    selector: { app: kafka }
    ports: [{port: 9092, targetPort: 9092}]
#####################################################################
# RABBITMQ
#####################################################################
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: rabbitmq
    namespace: proyecto2
  spec:
    replicas: 1
    selector: { matchLabels: { app: rabbitmq } }
    template:
      metadata: { labels: { app: rabbitmq } }
      spec:
        imagePullSecrets:
        - name: harbor-regcred 
        containers:
        - name: rabbitmq
          image: 34.55.73.46.nip.io/library/rabbitmq:latest
          ports:
            - {containerPort: 5672}   # AMQP
            - {containerPort: 15672}  # UI
          env:
            - {name: RABBITMQ_DEFAULT_USER, value: guest}
            - {name: RABBITMQ_DEFAULT_PASS, value: guest}
- apiVersion: v1
  kind: Service
  metadata:
    name: rabbitmq
    namespace: proyecto2
  spec:
    selector: { app: rabbitmq }
    ports:
      - {name: amqp, port: 5672,  targetPort: 5672}
      - {name: ui,   port: 15672, targetPort: 15672}
#####################################################################
# REDIS
#####################################################################
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: redis
    namespace: proyecto2
  spec:
    replicas: 1
    selector: { matchLabels: { app: redis } }
    template:
      metadata: { labels: { app: redis } }
      spec:
        imagePullSecrets:
        - name: harbor-regcred 
        containers:
        - name: redis
          image: 34.55.73.46.nip.io/library/redis:latest
          ports: [{containerPort: 6379}]
- apiVersion: v1
  kind: Service
  metadata:
    name: redis
    namespace: proyecto2
  spec:
    selector: { app: redis }
    ports: [{port: 6379, targetPort: 6379}]
#####################################################################
# VALKEY
#####################################################################
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: valkey
    namespace: proyecto2
  spec:
    replicas: 1
    selector: { matchLabels: { app: valkey } }
    template:
      metadata: { labels: { app: valkey } }
      spec:
        imagePullSecrets:
        - name: harbor-regcred 
        containers:
        - name: valkey
          image: 34.55.73.46.nip.io/library/valkey/valkey:latest
          ports: [{containerPort: 6379}]
- apiVersion: v1
  kind: Service
  metadata:
    name: valkey
    namespace: proyecto2
  spec:
    selector: { app: valkey }
    ports: [{port: 6379, targetPort: 6379}]
#####################################################################
# GRAFANA
#####################################################################
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: grafana
    namespace: proyecto2
  spec:
    replicas: 1
    selector: { matchLabels: { app: grafana } }
    template:
      metadata: { labels: { app: grafana } }
      spec:
        imagePullSecrets:
        - name: harbor-regcred 
        containers:
        - name: grafana
          image: 34.55.73.46.nip.io/library/grafana/grafana:latest
          ports: [{containerPort: 3000}]
- apiVersion: v1
  kind: Service
  metadata:
    name: grafana
    namespace: proyecto2
  spec:
    type: NodePort   # accesible fuera del cluster
    selector: { app: grafana }
    ports:
      - {port: 3000, targetPort: 3000, nodePort: 30000}
#####################################################################
# API-GO
#####################################################################
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: api-go
    namespace: proyecto2
  spec:
    replicas: 1
    selector: { matchLabels: { app: api-go } }
    template:
      metadata: { labels: { app: api-go } }
      spec:
        imagePullSecrets:
        - name: harbor-regcred 
        containers:
        - name: api-go
          image: 34.55.73.46.nip.io/library/proyecto2-api-go:latest
          ports: [{containerPort: 8081}]
- apiVersion: v1
  kind: Service
  metadata:
    name: api-go
    namespace: proyecto2
  spec:
    selector: { app: api-go }
    ports: [{port: 8081, targetPort: 8081}]
#####################################################################
# API-RUST  (depende de api-go →  readiness comprueba /health)
#####################################################################
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: api-rust
    namespace: proyecto2
  spec:
    replicas: 1
    selector: { matchLabels: { app: api-rust } }
    template:
      metadata: { labels: { app: api-rust } }
      spec:
        imagePullSecrets:
        - name: harbor-regcred 
        containers:
        - name: api-rust
          image: 34.55.73.46.nip.io/library/proyecto2-api-rust:latest
          ports: [{containerPort: 8082}]
          readinessProbe:   # espera a que api-go esté disponible
            httpGet:
              path: /health
              port: 8082
            initialDelaySeconds: 5
            periodSeconds: 5
- apiVersion: v1
  kind: Service
  metadata:
    name: api-rust
    namespace: proyecto2
  spec:
    selector: { app: api-rust }
    ports: [{port: 8082, targetPort: 8082}]
#####################################################################
# KAFKA WRITER
#####################################################################
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: kafka-writer
    namespace: proyecto2
  spec:
    replicas: 1
    selector: { matchLabels: { app: kafka-writer } }
    template:
      metadata: { labels: { app: kafka-writer } }
      spec:
        imagePullSecrets:
        - name: harbor-regcred 
        containers:
        - name: kafka-writer
          image: 34.55.73.46.nip.io/library/proyecto2-kafka-writer:latest
          ports: [{containerPort: 8083}]
- apiVersion: v1
  kind: Service
  metadata:
    name: kafka-writer
    namespace: proyecto2
  spec:
    selector: { app: kafka-writer }
    ports: [{port: 8083, targetPort: 8083}]
#####################################################################
# RABBIT WRITER
#####################################################################
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: rabbit-writer
    namespace: proyecto2
  spec:
    replicas: 1
    selector: { matchLabels: { app: rabbit-writer } }
    template:
      metadata: { labels: { app: rabbit-writer } }
      spec:
        imagePullSecrets:
        - name: harbor-regcred 
        containers:
        - name: rabbit-writer
          image: 34.55.73.46.nip.io/library/proyecto2-rabbit-writer:latest
          ports: [{containerPort: 8084}]
- apiVersion: v1
  kind: Service
  metadata:
    name: rabbit-writer
    namespace: proyecto2
  spec:
    selector: { app: rabbit-writer }
    ports: [{port: 8084, targetPort: 8084}]
#####################################################################
# KAFKA CONSUMER
#####################################################################
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: kafka-consumer
    namespace: proyecto2
  spec:
    replicas: 1
    selector: { matchLabels: { app: kafka-consumer } }
    template:
      metadata: { labels: { app: kafka-consumer } }
      spec:
        imagePullSecrets:
        - name: harbor-regcred 
        containers:
        - name: kafka-consumer
          image: 34.55.73.46.nip.io/library/proyecto2-kafka-consumer:latest
          ports: [{containerPort: 8085}]
- apiVersion: v1
  kind: Service
  metadata:
    name: kafka-consumer
    namespace: proyecto2
  spec:
    selector: { app: kafka-consumer }
    ports: [{port: 8085, targetPort: 8085}]
#####################################################################
# RABBIT CONSUMER
#####################################################################
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: rabbit-consumer
    namespace: proyecto2
  spec:
    replicas: 1
    selector: { matchLabels: { app: rabbit-consumer } }
    template:
      metadata: { labels: { app: rabbit-consumer } }
      spec:
        imagePullSecrets:
        - name: harbor-regcred 
        containers:
        - name: rabbit-consumer
          image: 34.55.73.46.nip.io/library/proyecto2-rabbit-consumer:latest
          ports: [{containerPort: 8086}]
- apiVersion: v1
  kind: Service
  metadata:
    name: rabbit-consumer
    namespace: proyecto2
  spec:
    selector: { app: rabbit-consumer }
    ports: [{port: 8086, targetPort: 8086}]
```

### Se hicieron 12 pods, con estos deployments, halando las imagenes directamente desde harbor, y ahi mismo se creó el servicio. haciendo replicas con ellos para realizar las pruebas correspondientes.

## ¿Cómo funciona Kafka?


### Apache Kafka es una plataforma de mensajería distribuida diseñada para manejar flujos de datos en tiempo real. Funciona como un sistema publicador suscriptor altamente escalable.

### Funcionando de la siguiente manera:

### - Productores envían mensajes a topics.

### - Los brokers almacenan esos mensajes de forma ordenada y persistente.

### - Los consumidores se suscriben a esos topics y leen mensajes en el orden en que fueron escritos.

### - Kafka asegura alta disponibilidad mediante replicación y permite procesamiento asíncrono y desacoplado entre servicios.

## ¿Cómo difiere Valkey de Redis?

### Valkey es un fork de Redis mantenido por la comunidad, creado en 2024 luego de que Redis pasó a ser un software bajo licencia comercial.

### Valkey es de codigo abierto y siguio el objetivo inicial de redis, ya que los colaboradores anteriores siguienron.

## ¿Es mejor gRPC que HTTP?

### Depende del uso, pero para servicios entre microservicios y comunicación interna, sí, gRPC es mejor que HTTP REST.

### Se usa gRPC en los servicios internos porque nos permite menor latencia, mayor rendimiento y tipado fuerte entre servicios.

##  ¿Hubo una mejora al utilizar dos replicas en los deployments de API REST y gRPC?

### Si hubo mejora, debido que con una replica con 10 usuarios cada 30 segundos hacia 30 peticiones aceptadas por locust, en cambio con las replicas llego a las 65 peticiones aceptadas por 30 segundos, y con mas de 100 usuarios se llegó a la meta.


## Para los consumidores, ¿Qué utilizó y por qué?

### primero se crearon los topics correspondientes si en dado caso no existenm, esto en la parte de los writers y el consumer solo escucha siempre cualquier cambio en el topic.

### - segmentio/kafka-go integración directa, alto rendimiento y confiabilidad al leer mensajes desde topics Kafka.

### - streadway/amqp facilidad para manejar colas durables y control de ack/manual requeue