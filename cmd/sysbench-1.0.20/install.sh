sysbench --db-driver=mysql --time=60 --threads=10 --report-interval=1 --mysql-host=192.168.15.69  --mysql-port=3306 --mysql-user=root --mysql-password=root123456 --mysql-db=tessan_erp --tables=20 --table_size=10000 oltp_read_write --db-ps-mode=disable prepare