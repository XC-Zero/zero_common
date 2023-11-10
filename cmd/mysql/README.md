# Mysql 主从配置

### 1. 第一步    启动 docker
   ```shell
    docker-compose up -d
   ```
### 2. 第二步 配置主节点
   1. 进入主节点容器
        ```shell
      docker exec -it 容器ID /bin/bash
      mysql -u root -p
        ```
   2. 新增同步账号
        ```mysql
      create user 'slave'@'%' identified with caching_sha2_password  by '123456';
      grant replication slave on *.* to 'slave'@'%';
      CHANGE MASTER TO GET_MASTER_PUBLIC_KEY=1;
      flush privileges;
        ```
   3. 重置 master 日志,记录 position
        ``` mysql
        Reset master;
        Show master status;
        Set global sql_mode="STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION"
        ```
### 3. 第三步 配置从节点
   1. 进入从节点容器 或以 root登录数据库
      ```shell
      docker exec -it 容器ID /bin/bash;
      mysql -u root -p
      ```
   2. 执行
      ```mysql
      SET GLOBAL server_id =  5321; # (不能和其他节点重复);
      Set global sql_mode="STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION";
      CHANGE MASTER TO GET_MASTER_PUBLIC_KEY=1;
      CHANGE MASTER TO MASTER_HOST='172.16.238.10(主节点IP)' , MASTER_PORT=3306 , MASTER_USER='slave' , MASTER_PASSWORD='123456' , MASTER_LOG_FILE='binlog.000001(主节点起始日志)' , MASTER_LOG_POS='29064(主节点 position)';
      ```
   3. 启动同步
      ```mysql
       start slave ;
       ```
### 4.第四步 查询状态
   1. 查看主节点状态
      ```mysql
      show master status ;
      ```
   2. 查看从节点状态
      ```mysql
      show slave status ;
      ```



### PS:  每次重启 从节点都需重新 设置 server_id 并启动同步
 
如遇报错   

`Host is blocked because of many connection errors`  

在主节点执行 
```mysql 
 flush hosts;
```
再重新启动主从