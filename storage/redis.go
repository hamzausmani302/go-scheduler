package storage

import (
	"context"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Db       int
}

type RedisStorage struct {
	config RedisConfig
	client *redis.Client
	ctx    *context.Context
	prefix string
}

// creates new instance of redis store
func NewRedisStorage(options RedisConfig) (postgres *RedisStorage, err error) {
	ctx := context.Background()
	redisStore := RedisStorage{config: options, ctx: &ctx, prefix: "scheduler"}
	if err := redisStore.connect(); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
		return nil, err
	}
	return &redisStore, nil
}

// Connect creates a redis connection and store it in client field.
func (redisStorage *RedisStorage) connect() (err error) {
	redisStorage.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisStorage.config.Host, redisStorage.config.Port), // e.g. "localhost:6379"
		Password: redisStorage.config.Password,                                             // "" if no password
		DB:       redisStorage.config.Db,                                                   // 0 is default DB
	})
	_, err1 := redisStorage.client.Ping(context.Background()).Result()
	if err1 != nil {
		log.Fatal(err1)
		return err1
	}
	return nil
}

func (redisStorage *RedisStorage) Close() error {
	if err := redisStorage.client.Close(); err != nil {
		return err
	}
	return nil
}

func (redisStorage *RedisStorage) Add(task TaskAttributes) error {
	//check if not already exist otherwise insert
	exists, err := redisStorage.client.Exists(*redisStorage.ctx, redisStorage.generateKey(task.Hash)).Result()
	if err != nil {
		return fmt.Errorf("error reading data from redis: %w", err)
	}
	if exists <= 0 {
		// key does not exist yet, insert it
		if err := redisStorage.insert(task); err != nil {
			return err
		}
	}
	return nil
}

func (redisStorage *RedisStorage) Fetch() ([]TaskAttributes, error) {
	// Lua script to make the queries faster as they will be executed by redis engine itself
	script := redis.NewScript(`
		local result = {}
		for i, key in ipairs(KEYS) do
			result[i] = redis.call("HGETALL", key)
		end
		return result
	`)
	keys, errk := redisStorage.getKeysByPattern(fmt.Sprintf("%s*", redisStorage.prefix))
	if errk != nil {
		log.Fatal("cannot get keys")
	}
	vals, err := script.Run(*redisStorage.ctx, redisStorage.client, keys).Result()
	if err != nil {
		log.Fatal(err)
		return []TaskAttributes{}, err
	}
	rawList, ok := vals.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", vals)
	}
	var tasks []TaskAttributes

	for _, raw := range rawList {
		arr, ok := raw.([]interface{})
		if !ok {
			continue
		}

		m := make(map[string]string)
		for i := 0; i < len(arr); i += 2 {
			field, _ := arr[i].(string)
			value, _ := arr[i+1].(string)
			m[field] = value
		}

		task := TaskAttributes{
			Hash:        m["hash"],
			Name:        m["name"],
			LastRun:     m["lastrun"],
			NextRun:     m["nextrun"],
			Duration:    m["duration"],
			IsRecurring: m["isrecurring"],
			Params:      m["params"],
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (redisStorage *RedisStorage) Remove(task TaskAttributes) error {
	return redisStorage.client.Del(*redisStorage.ctx, redisStorage.generateKey(task.Hash)).Err()
}

func (redisStorage *RedisStorage) insert(task TaskAttributes) (err error) {
	bytesData, marshallErr := yaml.Marshal(task)
	if marshallErr != nil {
		return marshallErr
	}
	result := map[string]string{}
	if err := yaml.Unmarshal(bytesData, result); err != nil {
		return err
	}
	return redisStorage.client.HMSet(*redisStorage.ctx, redisStorage.generateKey(task.Hash), result).Err()
}

func (redisStorage *RedisStorage) getKeysByPattern(pattern string) ([]string, error) {
	var (
		cursor uint64
		keys   []string
	)
	for {
		k, nextCursor, err := redisStorage.client.Scan(*redisStorage.ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, k...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

func (redisStorage *RedisStorage) generateKey(hash string) string {
	return fmt.Sprintf("%s-%s", redisStorage.prefix, hash)
}
