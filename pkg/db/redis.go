package db

import (
	"context"
	"log"

	"tron-scan/config"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis(cfg config.RedisConfig) {
	RDB = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	if _, err := RDB.Ping(Ctx).Result(); err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}

	log.Println("Redis连接成功")
}
