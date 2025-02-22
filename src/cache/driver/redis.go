package driver

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/trancecho/open-sdk/cache/types"
	"github.com/trancecho/open-sdk/config"
	"github.com/trancecho/open-sdk/model"
	"log"
	"time"
)

type RedisCreator struct{}

func (c RedisCreator) Create(conf config.Cache) (types.Cache, error) {
	var r RedisCache
	r.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.IP, conf.PORT),
		Password: conf.PASSWORD,
		DB:       conf.DB,
	})
	_, err := r.Client.Ping().Result()
	if err != nil {
		log.Fatalln(err)
	}
	return r, nil
}

type RedisCache struct {
	Client *redis.Client
}

func (r RedisCache) GetInt(key string) (int, bool) {
	value, err := r.Client.Get(key).Int()
	if err == nil {
		return value, true
	}
	if err != redis.Nil {
		log.Fatalln(err)
	}

	return 0, false
}

func (r RedisCache) GetInt64(key string) (int64, bool) {
	value, err := r.Client.Get(key).Int64()
	if err == nil {
		return value, true
	}
	if err != redis.Nil {
		log.Fatalln(err)
	}
	return 0, false
}

func (r RedisCache) GetFloat32(key string) (float32, bool) {
	value, err := r.Client.Get(key).Float32()
	if err == nil {
		return value, true
	}
	if err != redis.Nil {
		log.Fatalln(err)
	}
	return 0, false
}

func (r RedisCache) GetFloat64(key string) (float64, bool) {
	value, err := r.Client.Get(key).Float64()
	if err == nil {
		return value, true
	}
	if err != redis.Nil {
		log.Fatalln(err)
	}
	return 0, false
}

func (r RedisCache) GetString(key string) (string, bool) {
	value, err := r.Client.Get(key).Result()
	if err == nil {
		return value, true
	}
	if err != redis.Nil {
		log.Fatalln(err)
	}
	return "", false
}

func (r RedisCache) GetBool(key string) (bool, bool) {
	value, err := r.Client.Get(key).Result()
	if err != redis.Nil {
		log.Fatalln(err)
	}
	if value == "1" {
		return true, true
	} else if value == "0" {
		return false, true
	}
	return false, false
}

func (r RedisCache) Set(Key string, value any, expireDuration time.Duration) error {
	return r.Client.Set(Key, value, expireDuration).Err()
}

func (r RedisCache) Del(key string) bool {
	err := r.Client.Del(key).Err()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Fatalln(err)
	}
	return true
}

func (r RedisCache) GetRedis() *redis.Client {
	return r.Client
}

func (r RedisCache) Exists(key string) bool {
	_, err := r.Client.Get(key).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Fatalln(err)
	}
	return true
}

// InsertMessage 插入消息
func InsertMessage(rdb *redis.Client, id string, content string, read bool, from string) error {
	key := fmt.Sprintf("message:%s", id)
	comment := &model.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		Read:      read,
		From:      from,
	}
	commentJSON, err := json.Marshal(comment)
	if err != nil {
		log.Printf("json marshal failed: %v", err)
		return err
	}
	// 使用ZADD存储消息
	err = rdb.ZAdd(key, redis.Z{
		Score:  float64(comment.StartTime),
		Member: commentJSON,
	}).Err()
	if err != nil {
		log.Printf("zadd failed: %v", err)
		return err
	}
	return nil
}

// UpdateAllMessageReadStatus 更新数据库中消息的读取状态，先删除所有元素，然后重新添加它们并更新已读状态(reserved)
func UpdateAllMessageReadStatus(rdb *redis.Client, id string) error {
	key := fmt.Sprintf("message:%s", id)
	messages, err := FirstGetMany(rdb, id)
	if err != nil {
		return err
	}

	// 删除集合中的所有元素
	err = rdb.Del(key).Err()
	if err != nil {
		log.Printf("Failed to delete all elements: %v", err)
		return err
	}

	// 更新读取状态并重新添加消息
	for _, message := range messages {
		message.Read = true
		messageJSON, err2 := json.Marshal(message)
		if err2 != nil {
			log.Printf("json marshal failed: %v", err)
			return err
		}
		err = rdb.ZAdd(key, redis.Z{
			Score:  float64(message.StartTime),
			Member: messageJSON,
		}).Err()
		if err != nil {
			log.Printf("zadd failed: %v", err)
			return err
		}
	}
	return nil
}

// FirstGetMany GetMessage 获取消息
func FirstGetMany(rdb *redis.Client, ID string) ([]model.Trainer, error) {
	// 获取所有的消息
	keyMe := fmt.Sprintf("message:%s", ID)
	now := time.Now().Unix()
	results, err := rdb.ZRangeByScore(keyMe, redis.ZRangeBy{
		Min: "-inf",                 //使用负无穷大为最小值
		Max: fmt.Sprintf("%d", now), //现在
	}).Result()
	if err != nil {
		return nil, err
	}

	var messages []model.Trainer
	for _, result := range results {
		var message model.Trainer
		if err = json.Unmarshal([]byte(result), &message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// DeleteMessages 删除过期消息
func DeleteMessages(rdb *redis.Client, ID string) error {
	key := fmt.Sprintf("message:%s", ID)
	now := time.Now().Unix()
	err := rdb.ZRemRangeByScore(key, "-inf", fmt.Sprintf("%d", now)).Err()
	if err != nil {
		log.Printf("Failed to delete expired messages: %v", err)
		return err
	}
	log.Printf("Delete %s messages successfully", key)
	return nil
}
