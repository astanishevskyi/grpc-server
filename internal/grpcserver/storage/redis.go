package storage

import (
	"context"
	"encoding/json"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/models"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
)

type RedisStorage struct {
	RedisStorage *redis.Client
	lastID       uint32
}

func NewRedisStorage() *RedisStorage {
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASS")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	storage := &RedisStorage{RedisStorage: rdb}
	lastID := storage.getLastID()
	storage.lastID = lastID
	return storage
}

func (r *RedisStorage) getLastID() uint32 {
	var cursor uint64
	var lastID uint32
	var err error
	var result []string
	for {
		result, cursor, err = r.RedisStorage.Scan(context.Background(), cursor, "*", 10).Result()
		if err != nil {
			log.Fatal(err)
		}
		for _, i := range result {
			id, err := strconv.Atoi(i)
			if err != nil {
				log.Fatal(err)
			}
			if lastID < uint32(id) {
				lastID = uint32(id)
			}
		}
		if cursor == 0 {
			break
		}
	}
	return lastID
}

func (r *RedisStorage) GetAll() ([]models.User, error) {
	var cursor uint64
	var result []models.User
	for {
		res, cursor, err := r.RedisStorage.Scan(context.Background(), cursor, "*", 10).Result()
		for _, i := range res {
			i, err := strconv.Atoi(i)
			if err != nil {
				return nil, err
			}
			retrieve, err := r.Retrieve(uint32(i))
			if err != nil {
				return nil, err
			}
			result = append(result, retrieve)
		}
		if err != nil {
			return nil, err
		}
		if cursor == 0 {
			break
		}
	}

	return result, nil
}

func (r *RedisStorage) Retrieve(id uint32) (models.User, error) {
	val, err := r.RedisStorage.Get(context.Background(), strconv.Itoa(int(id))).Result()
	if err != nil {
		return models.User{}, err
	}
	user := models.User{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *RedisStorage) Add(name, email string, age uint8) (models.User, error) {
	r.lastID++
	newUser := models.User{
		ID:    r.lastID,
		Name:  name,
		Email: email,
		Age:   age,
	}
	jsonUser, err := json.Marshal(newUser)
	if err != nil {
		return models.User{}, err
	}
	err = r.RedisStorage.Set(context.Background(), strconv.Itoa(int(r.lastID)), jsonUser, 0).Err()
	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

func (r *RedisStorage) Remove(id uint32) (uint32, error) {
	err := r.RedisStorage.Get(context.Background(), strconv.Itoa(int(id))).Err()
	if err != nil {
		return 0, err
	}
	err = r.RedisStorage.Del(context.Background(), strconv.Itoa(int(id))).Err()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *RedisStorage) Update(id uint32, name, email string, age uint8) (models.User, error) {
	err := r.RedisStorage.Get(context.Background(), strconv.Itoa(int(id))).Err()
	if err != nil {
		return models.User{}, err
	}
	newUser := models.User{
		ID:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}
	jsonUser, err := json.Marshal(newUser)
	if err != nil {
		return models.User{}, err
	}
	err = r.RedisStorage.Set(context.Background(), strconv.Itoa(int(id)), jsonUser, 0).Err()
	if err != nil {
		return models.User{}, err
	}
	r.lastID++
	return newUser, nil
}
