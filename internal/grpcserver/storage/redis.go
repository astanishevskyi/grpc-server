package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/models"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type RedisStorage struct {
	Addr         string
	Password     string
	DB           int
	RedisStorage *redis.Client
	lastID       uint32
}

func NewRedisStorage() *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// get last id
	return &RedisStorage{RedisStorage: rdb, lastID: 1}
}

func (r *RedisStorage) GetAll() []models.User {
	// iterate using scan
	return nil
}

func (r *RedisStorage) Retrieve(id uint32) (models.User, error) {
	val, err := r.RedisStorage.Get(context.Background(), strconv.Itoa(int(id))).Result()
	if err != nil {
		return models.User{}, err
	}
	fmt.Println(val)
	user := models.User{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *RedisStorage) Add(name, email string, age uint8) (models.User, error) {
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
	r.lastID++

	return newUser, nil
}

func (r *RedisStorage) Remove(id uint32) (uint32, error) {
	err := r.RedisStorage.Del(context.Background(), strconv.Itoa(int(id))).Err()
	if err != nil {
		return 0, err
	}
	// make scan
	return id, nil
}

func (r *RedisStorage) Update(id uint32, name, email string, age uint8) (models.User, error) {
	//make scan
	_, err := r.RedisStorage.Get(context.Background(), strconv.Itoa(int(id))).Result()
	if err != nil {
		return models.User{}, err
	}
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
	r.lastID++
	return newUser, nil
}
