package models

//redis官网: github.com/go-redis
//下载go-redis: go get github.com/redis/go-redis/v9
//连接redis数据库核心代码

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// 全局使用,就需要把定义成公有的
var ctxRedis = context.Background()

var (
	RedisDb *redis.Client
)

// 自动初始化数据库
func init() {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//连接redis
	_, err := RedisDb.Ping(ctxRedis).Result()
	//判断连接是否成功
	if err != nil {
		println(err)
	}
}
