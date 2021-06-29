package storage

import (
	"context"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoStorage struct {
	Addr         string
	Password     string
	DB           int
	MongoStorage *mongo.Client
}

func NewMongoStorage() *MongoStorage {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	return &MongoStorage{MongoStorage: client}
}

func (r *MongoStorage) GetAll() []models.User {
	return nil
}

func (r *MongoStorage) Retrieve(id uint32) (models.User, error) {
	return models.User{}, nil
}

func (r *MongoStorage) Add(name, email string, age uint8) (models.User, error) {
	return models.User{}, nil
}

func (r *MongoStorage) Remove(uint32) (uint32, error) {
	return 0, nil
}

func (r *MongoStorage) Update(id uint32, name, email string, age uint8) (models.User, error) {
	return models.User{}, nil
}
