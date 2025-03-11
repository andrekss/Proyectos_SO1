#!/bin/bash
#../Scripts/Contenedores.sh  # comentar si se usa en rust

docker ps -a
# ------ comandos específico para rust ------
cd .. 
cd module
# ------ comandos específico para rust ------
make
sudo insmod sysinfo_202113580.ko
sudo dmesg | tail
cat /proc/sysinfo_202113580
sudo rmmod sysinfo_202113580
make clean