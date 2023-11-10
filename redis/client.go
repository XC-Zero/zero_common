package redis

import (
	"github.com/go-redis/redis/v7"
	"gitlab.tessan.com/data-center/tessan-erp-common/config"
)

// InitRedis ...
func InitRedis(config config.RedisConfig) (*redis.Client, error) {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.AddrList[0],
		Password: config.Password,
		DB:       0,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	return redisClient, nil
}

// InitRedisCluster 初始化Redis集群的连接
func InitRedisCluster(config config.RedisConfig) (*redis.ClusterClient, error) {

	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    config.AddrList,
		Password: config.Password,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	return redisClient, nil
}
