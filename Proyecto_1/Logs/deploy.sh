#!/bin/bash

Nombre_Imagen="dockerfile"  
Nombre_Contenedor="Servicio"  

echo "Deteniendo y eliminando contenedor existente (si aplica)..."
docker stop $Nombre_Contenedor 2>/dev/null
docker rm $Nombre_Contenedor 2>/dev/null

echo "Eliminando imagen anterior (si existe)..."
docker rmi -f $Nombre_Imagen 2>/dev/null

echo "Construyendo la nueva imagen..."
docker build -t $Nombre_Imagen .

echo "Corriendo el nuevo contenedor..."
docker run -d -p 8000:8000 --name $Nombre_Contenedor $Nombre_Imagen

echo "Despliegue completado. El contenedor está corriendo con la nueva imagen."
