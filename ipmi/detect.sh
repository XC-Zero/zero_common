
#查看主板型号
dmidecode -t baseboard | grep "Product Name"
#查看CPU型号
cat /proc/cpuinfo | grep "model name" | uniq
#查看IPMI情况
dmidecode |grep -i ipmi

#尝试开启IPMI
sudo modprobe ipmi_watchdog
sudo modprobe ipmi_poweroff
sudo modprobe ipmi_devintf