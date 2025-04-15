
source entorno/bin/activate
locust -f locustfile.py --host=http://192.168.238.35.nip.io

## deactivate


# Instalar ngress Controller (NGINX)
# helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
# helm repo update
#helm install ingress-nginx ingress-nginx/ingress-nginx \
#  --create-namespace \
#  --namespace ingress-nginx


# kubectl get svc -n ingress-nginx // verificar que corre
# kubectl get pod -n ingress-nginx // verificar pods

# kubectl get ingress