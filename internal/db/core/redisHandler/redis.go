package redisHandler

import (
	"context"
	"fmt"
	"ipfast_server/pkg/util/log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/spf13/viper"
)

var Rdb *redis.Client
var ctx = context.Background()

func init() {
	Rdb = &redis.Client{}
}

/*
初始化redis连接
*/
func Setup() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.host"),     // Redis 地址
		Password:     viper.GetString("redis.password"), // 密码，没有则留空
		DB:           viper.GetInt("redis.db"),          // 使用默认DB
		PoolSize:     viper.GetInt("redis.poolSize"),    // 连接池最大连接数
		MinIdleConns: 5,                                 // 在连接池中维护的最小空闲连接数
		PoolTimeout:  30 * time.Second,                  // 客户端等待空闲连接的最长时间
	})
	code := Rdb.Ping(ctx)
	if code.Err() != nil {
		log.Fatalln("redis连接失败:%v", code.Err())
	}
	log.Info("redis连接成功")
}

// 使用事务和乐观锁 累加更新值，并加入重试机制
func UpdateWithAccumulation(key string, value int64, expiration time.Duration) (err error) {
	const maxRetries = 5 // 最大重试次数
	retryCount := 0      // 当前重试次数

	for {
		err = Rdb.Watch(ctx, func(tx *redis.Tx) error {
			n, err := tx.Get(ctx, key).Int64()
			if err != nil && err != redis.Nil {
				return err
			}
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				n += value
				pipe.Set(ctx, key, n, expiration)
				return nil
			})
			return err
		}, key)

		if err == nil {
			return
		}

		if err == redis.TxFailedErr {
			retryCount++
			log.Error("由于乐观锁失败，更新重试: %d", retryCount)
			if retryCount >= maxRetries {
				return fmt.Errorf("达到最大重试次数[%d]，更新失败[%s]:%d", maxRetries, key, value)
			}
		} else {
			return fmt.Errorf("更新失败[%s]: %v", key, err)
		}
	}
}

func SetInt64(key string, value int64, expiration time.Duration) error {
	err := Rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Error("更新失败[%s]: %v", key, err)
		return err
	}
	return nil
}

func GetInt64(key string) (value int64, err error) {
	value, err = Rdb.Get(ctx, key).Int64()
	if err != nil {
		return
	}
	return value, nil
}

func SetInt(key string, value int, expiration time.Duration) error {
	err := Rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Error("更新失败[%s]: %v", key, err)
		return err
	}
	return nil
}

func GetInt(key string) (value int, err error) {
	value, err = Rdb.Get(ctx, key).Int()
	if err != nil {
		return
	}
	return value, nil
}

func HSet(key string, mapData map[string]interface{}) error {
	redisData := make([]interface{}, 0, len(mapData)*2)
	for k, v := range mapData {
		redisData = append(redisData, k, v)
	}
	err := Rdb.HSet(ctx, key, redisData...).Err()
	if err != nil {
		return err
	}
	return nil
}

func HGet(key string) (data map[string]string, err error) {
	data, err = Rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return
	}
	return
}
