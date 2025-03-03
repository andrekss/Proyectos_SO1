crontab -e # edicion
../Scripts/Contenedores.sh
docker ps -a
sleep 5
make
sudo insmod sysinfo_202113580.ko
sudo dmesg | tail
cat /proc/sysinfo_202113580
sudo rmmod sysinfo_202113580