package storage

//
//import (
//	"context"
//	"github.com/astanishevskyi/grpc-server/internal/grpcserver/models"
//	"github.com/go-redis/redis/v8"
//)
//
//type MongoStorage struct {
//	Addr string
//	Password string
//	DB int
//	MongoStorage *redis.Client
//}
//
//func NewMongoStorage() *MongoStorage {
//	rdb := *redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "", // no password set
//		DB:       0,  // use default DB
//	})
//	err := rdb.Set(context.Background(), "key", "value", 0).Err()
//	if err != nil {
//		panic(err)
//	}
//	return &MongoStorage{MongoStorage: &rdb}
//
//}
//
//func (r *MongoStorage) GetAll() []models.User {
//	return nil
//}
//
//func (r *MongoStorage) Retrieve(id uint32) (models.User, error) {
//	return models.User{}, nil
//}
//
//func (r *MongoStorage) Add(name, email string, age uint8) (models.User, error) {
//	return models.User{}, nil
//}
//
//func (r *MongoStorage) Remove(uint32) (uint32, error) {
//	return 0, nil
//}
//
//func (r *MongoStorage) Update(id uint32, name, email string, age uint8) (models.User, error) {
//	return models.User{}, nil
//}
//
