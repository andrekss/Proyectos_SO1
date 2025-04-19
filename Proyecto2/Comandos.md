gcloud auth login
gcloud container clusters list

// Conectar al clúster
gcloud container clusters get-credentials so-cluster-1 \
  --zone us-central1-c \
  --project reliable-byte-455302-k8


sudo du -h / --max-depth=2 | sort -hr | head -20 // limpieza
docker-compose up --build // subir y construir