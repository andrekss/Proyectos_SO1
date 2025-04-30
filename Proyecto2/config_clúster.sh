kubectl create namespace proyecto2

# Autenticación de harbor
kubectl create secret docker-registry harbor-regcred \
  --docker-server=34.55.73.46.nip.io \
  --docker-username=admin \
  --docker-password=Harbor12345 \
  --docker-email=alejandroagosto11003@gmail.com \
  -n proyecto2

# todo el despliegue
kubectl apply -f deployment.yml

# verifique pods y servicios
kubectl get pods -n proyecto2
kubectl get svc  -n proyecto2

# inspeccionar logs
kubectl logs -f deployment/kafka -n proyecto2


# ingress controller nip.io 

# Agregar el repositorio de NGINX Ingress Controller
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx

helm repo update

kubectl create namespace ingress-nginx
helm install ingress-nginx ingress-nginx/ingress-nginx --namespace ingress-nginx


kubectl apply -f ingress.yml

# verificar existencia
kubectl get ingress -n proyecto2

# http://34.118.236.112.nip.io/input



# extra

kubectl describe pod <pod-name> -n proyecto2

# modificar réplicas
kubectl scale deployment api-go --replicas=1 -n proyecto2

# jalar crt

openssl s_client -showcerts -connect 34.55.73.46:443 </dev/null 2>/dev/null | openssl x509 -outform PEM > harbor.crt
sudo mkdir -p /etc/docker/certs.d/34.55.73.46
sudo cp harbor.crt /etc/docker/certs.d/34.55.73.46/

# borrar pods
kubectl delete pods --all -n proyecto2

# borrar todos los pods
kubectl delete pods --all -n proyecto2

# exposicion de ingress y certificados en el siguiente video https://drive.google.com/file/d/1cn71C6DHKgFvyPrzPPkiGx0n2gqI4YiF/view?t=3755

# Conclusiones: los certificados necesitan de ip.nip.io para evitar errores usar certbot
# Para hacer el ingress la ip que se usa es la del ingress-controller