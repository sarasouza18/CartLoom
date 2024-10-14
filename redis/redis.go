package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// InitRedis initializes a connection to Redis
func InitRedis(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Address of the Redis server
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis")
	return rdb, nil
}

// ProcessOrder processes an order and updates Redis with a lock mechanism
func ProcessOrder(ctx context.Context, rdb *redis.Client, orderID string) error {
	lockKey := fmt.Sprintf("lock:%s", orderID)

	if acquired, err := acquireLock(ctx, rdb, lockKey); err != nil {
		return err
	} else if !acquired {
		return fmt.Errorf("could not acquire lock for order %s", orderID)
	}
	defer releaseLock(ctx, rdb, lockKey)

	if err := updateOrderStatus(ctx, rdb, orderID, "Order Processed"); err != nil {
		return err
	}

	log.Printf("Order %s processed and status updated in Redis", orderID)
	return nil
}

// acquireLock tries to set a lock with a TTL to ensure the order is not processed concurrently
func acquireLock(ctx context.Context, rdb *redis.Client, lockKey string) (bool, error) {
	success, err := rdb.SetNX(ctx, lockKey, "locked", 10*time.Second).Result()
	if err != nil {
		return false, fmt.Errorf("error acquiring lock for %s: %v", lockKey, err)
	}
	return success, nil
}

// releaseLock removes the lock after the order is processed
func releaseLock(ctx context.Context, rdb *redis.Client, lockKey string) {
	if err := rdb.Del(ctx, lockKey).Err(); err != nil {
		log.Printf("Failed to release lock for %s: %v", lockKey, err)
	}
}

// updateOrderStatus updates the order status in Redis
func updateOrderStatus(ctx context.Context, rdb *redis.Client, orderID, status string) error {
	if err := rdb.Set(ctx, orderID, status, 0).Err(); err != nil {
		return fmt.Errorf("failed to update order %s in Redis: %v", orderID, err)
	}
	return nil
}

// InitRedisWithAddress initializes a Redis connection using a custom address
func InitRedisWithAddress(ctx context.Context, redisAddress string) (*redis.Client, error) {
	// Create a new Redis client using the address from the environment variable
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddress, // Address passed as a parameter
	})

	// Ping the Redis server to verify the connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis at %s: %v", redisAddress, err)
	}

	log.Println("Successfully connected to Redis at", redisAddress)
	return rdb, nil
}
