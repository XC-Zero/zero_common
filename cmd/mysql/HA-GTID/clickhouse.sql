SET allow_experimental_database_materialize_mysql=1;
CREATE DATABASE slave_db ENGINE = MaterializeMySQL('172.165.212.109:53306', 'test01x', 'copy_root', 'copy_root');
