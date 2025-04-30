# ---------------------- Conexion y versiones de docker ----------------------

# Conéctate a la VM
gcloud compute ssh harbor-vm

# Actualiza el sistema
sudo apt update && sudo apt upgrade -y

# Instala Docker
sudo apt install -y docker.io
sudo systemctl enable docker
sudo systemctl start docker

# Agrega tu usuario al grupo docker (requiere reconexión)
sudo usermod -aG docker $USER

# Instala Docker Compose
sudo apt install -y docker-compose



# ---------------------- Despliegue de harbor ----------------------

# Certificados Https
sudo mkdir -p /etc/harbor/ssl
sudo openssl req -newkey rsa:2048 -nodes -sha256 \
  -keyout /etc/harbor/ssl/harbor.key \
  -x509 -days 365 \
  -out /etc/harbor/ssl/harbor.crt \
  -subj "/C=GT/ST=Guatemala/L=Ciudad/O=Universidad/CN=harbor.local"

# Clona el repositorio oficial
git clone https://github.com/goharbor/harbor.git
cd harbor

# Descarga el instalador
wget https://github.com/goharbor/harbor/releases/download/v2.10.0/harbor-online-installer-v2.10.0.tgz
tar -xvzf harbor-online-installer-v2.10.0.tgz
cd harbor

# Copia el archivo de ejemplo
cp harbor.yml.tmpl harbor.yml

# Edita el archivo harbor.yml 
nano harbor.yml

# Escribir esto en el yml
#hostname: 34.xx.xx.xx
#http:
#  port: 80
#https:
#  port: 443
#  certificate: /etc/harbor/ssl/harbor.crt
#  private_key: /etc/harbor/ssl/harbor.key

# Ejecuta el script de instalación
sudo ./install.sh


# En la maquina local ejecutar el siguiente comando 
sudo nano /etc/docker/daemon.json

# Escribir lo siguiente
#{
#  "insecure-registries" : ["34.55.73.46"] # o la ip de la vm
#}
sudo systemctl restart docker
docker login 34.55.73.46 # ip de la vm para conectar en la local





# user: admin password: Harbor12345
