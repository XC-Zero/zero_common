
set global server_id =65431;
CREATE USER if not exists 'copy_root' IDENTIFIED BY 'copy_root';
grant all on *.* to 'copy_root'@'%';
stop slave;

CHANGE MASTER TO GET_MASTER_PUBLIC_KEY =1;
flush privileges;
reset master;
show master status;


change master to master_host ='172.16.212.102', master_user ='copy_root',master_password ='copy_root',master_port =3306 ,MASTER_LOG_FILE ='MYSQL_HA_MASTER_02-bin.000001' , MASTER_LOG_POS =157;
start slave;
show slave status;

