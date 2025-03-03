# Proyectos_SO1
### Crear entorno arch para pip
#### python -m venv entorno      // crear entorno
#### source entorno/bin/activate // entrar al entorno
#### deactivate

### Comandos modulo Kernel

#### cd module
#### make
#### sudo insmod sysinfo_202113580.ko
#### sudo dmesg | tail
#### cat /proc/sysinfo_202113580 // visualizar
#### sudo rmmod sysinfo_202113580 // desinstalar

### Logs
#### sudo docker build -t dockerfile . // generar imagen
#### docker run -d -p 8000:8000 dockerfile // ejecutar imagen

#### docker start id_conteiner
#### docker stop id_conteiner
