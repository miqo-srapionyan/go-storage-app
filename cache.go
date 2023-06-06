package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

// Passing redisClient with argument, so it can be mocked for testing

func getPromotionFromCache(id string, redisClient *redis.Client) (*Promotion, error) {
	cacheKey := redisCacheKeyPrefix + id

	// Check if the promotion exists in cache
	exists, err := redisClient.Exists(context.Background(), cacheKey).Result()
	if err != nil {
		return nil, err
	}

	if exists == 0 {
		return nil, redis.Nil
	}

	// Retrieve the promotion from cache
	promotionData, err := redisClient.HGetAll(context.Background(), cacheKey).Result()
	if err != nil {
		return nil, err
	}

	// Parse the promotion data
	price, err := strconv.ParseFloat(promotionData["Price"], 64)
	if err != nil {
		return nil, err
	}

	// Update the expiration time for the promotion in the sorted set
	redisClient.ZAdd(context.Background(), redisCacheKeyPrefix+"timestamps", redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: cacheKey,
	})

	return &Promotion{
		UUID:           promotionData["UUID"],
		Price:          price,
		ExpirationDate: promotionData["ExpirationDate"],
	}, nil
}

func addPromotionToCache(id string, promotion *Promotion) error {
	cacheKey := redisCacheKeyPrefix + id

	// Add the promotion to the cache hash map
	promotionData := map[string]interface{}{
		"UUID":           promotion.UUID,
		"Price":          promotion.Price,
		"ExpirationDate": promotion.ExpirationDate,
	}

	err := redisClient.HMSet(context.Background(), cacheKey, promotionData).Err()
	if err != nil {
		return err
	}

	err = redisClient.Expire(context.Background(), cacheKey, cacheExpiration).Err()
	if err != nil {
		return err
	}

	// Update the expiration time for the promotion in the sorted set
	redisClient.ZAdd(context.Background(), redisCacheKeyPrefix+"timestamps", redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: cacheKey,
	})

	// Check if the number of cache keys exceeds the cache size
	count, err := redisClient.ZCard(context.Background(), redisCacheKeyPrefix+"timestamps").Result()
	if err != nil {
		return err
	}

	if count > cacheSize {
		// Get the oldest cache key
		oldestKeys, err := redisClient.ZRange(context.Background(), redisCacheKeyPrefix+"timestamps", 0, 0).Result()
		if err != nil {
			return err
		}

		if len(oldestKeys) > 0 {
			// Delete the oldest promotion from the cache
			err = redisClient.Del(context.Background(), oldestKeys[0]).Err()
			if err != nil {
				return err
			}

			// Delete the corresponding timestamp from the sorted set
			err = redisClient.ZRem(context.Background(), redisCacheKeyPrefix+"timestamps", oldestKeys[0]).Err()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
