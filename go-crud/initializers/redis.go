// initializers/redis.go
package initializers

import (
    "context"
    "log"
    "time"

    "github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // No password
        DB:       0,  // Default DB
    })

    ctx := context.Background()
    pong, err := RedisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatal("Error connecting to Redis:", err)
    }
    log.Println("Connected to Redis:", pong)
}


func GetFromCache(key string) (string, error) {
    ctx := context.Background()
    return RedisClient.Get(ctx, key).Result()
}

func SetToCache(key string, value string, expiration time.Duration) error {
    ctx := context.Background()
    return RedisClient.Set(ctx, key, value, expiration).Err()
}
