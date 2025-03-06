../Scripts/Contenedores.sh
docker ps -a
make
sudo insmod sysinfo_202113580.ko
cat /proc/sysinfo_202113580
sudo rmmod sysinfo_202113580