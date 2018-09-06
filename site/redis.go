package main

import (
    "github.com/go-redis/redis"
    "time"
)

func getRedisConnection() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     config.RedisConfig.Host,
        Password: config.RedisConfig.Password,
        DB:       config.RedisConfig.Database,
    })
    return client
}

func (client *Redis) getRedis(key string) string {
    if val, err := client.Get(key).Result(); err != nil {
        panic(err)
    } else {
        return val
    }
}

func (client *Redis) setRedis(key string, val string) {
    client.setRedisTTL(key, val, 0)
}

func (client *Redis) setRedisTTL(key string, val string, ttl int) {
    if err := client.Set(key, val, time.Duration(ttl) * time.Second); err != nil {
        panic(err)
    }
}

func (client *Redis) setRedisPTTL(key string, val string, pttl int) {
    if err := client.Set(key, val, time.Duration(pttl) * time.Millisecond); err != nil {
        panic(err)
    }
}