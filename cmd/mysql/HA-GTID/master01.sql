
CREATE USER if not exists 'copy_root' IDENTIFIED WITH mysql_native_password  BY 'copy_root';
grant all on *.* to 'copy_root'@'%';
stop slave;
CHANGE MASTER TO GET_MASTER_PUBLIC_KEY =1;
 flush privileges;
reset master;
show master status;
stop slave;

change master to master_host='172.165.212.102',master_port=3306,MASTER_USER='copy_root',MASTER_PASSWORD='copy_root',MASTER_AUTO_POSITION=1;

start slave;
show slave status;

