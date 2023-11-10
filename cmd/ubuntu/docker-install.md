##   

## 基础工具类

```shell

sudo apt install openssh-server
sudo service ssh start
sudo apt install net-tools

```

## 硬盘类

### 1.安装duf

### 2. 检查硬盘挂载

```shell
duf
```

```shell
parted -ll
fdisk -l
```

### 如需初始化

```mkfs.xfs -f /dev/sdb```
