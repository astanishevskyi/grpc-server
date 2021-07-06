package storage

import (
	"context"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"sync"
)

type MongoStorage struct {
	Addr         string
	Password     string
	MongoStorage *mongo.Collection
	mu           sync.Mutex
	lastID       uint32
}

func NewMongoStorage() *MongoStorage {
	addr := os.Getenv("MONGO_ADDR")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://"+addr))
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("grpc_db").Collection("users")
	storage := &MongoStorage{MongoStorage: collection}
	id, err := storage.getLastID()
	if err != nil {
		log.Fatal(err)
	}
	storage.lastID = id
	return storage
}

func (m *MongoStorage) getLastID() (uint32, error) {
	var result bson.M
	opts := options.FindOne().SetSort(bson.D{{"_id", -1}})
	err := m.MongoStorage.FindOne(context.Background(), bson.D{}, opts).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}
	id := uint32(result["_id"].(int64))
	return id, nil
}

func (m *MongoStorage) GetAll() ([]models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	findUsers, err := m.MongoStorage.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = findUsers.All(context.Background(), &results); err != nil {
		return nil, err
	}
	var users []models.User
	for _, result := range results {
		_id := uint32(result["_id"].(int64))
		name := result["name"].(string)
		age := uint8(result["age"].(int32))
		email := result["email"].(string)

		user := models.User{ID: _id, Name: name, Age: age, Email: email}
		users = append(users, user)
	}
	return users, nil
}

func (m *MongoStorage) Retrieve(id uint32) (models.User, error) {
	var result bson.M
	err := m.MongoStorage.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return models.User{}, err
	}
	_id := uint32(result["_id"].(int32))
	name := result["name"].(string)
	age := uint8(result["age"].(int32))
	email := result["email"].(string)

	user := models.User{ID: _id, Name: name, Age: age, Email: email}
	return user, nil
}

func (m *MongoStorage) Add(name, email string, age uint8) (models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	newUser := models.User{Name: name, Email: email, Age: age}
	m.lastID++
	res, err := m.MongoStorage.InsertOne(
		context.Background(),
		bson.M{
			"_id":   m.lastID,
			"name":  newUser.Name,
			"email": newUser.Email,
			"age":   newUser.Age,
		},
	)
	if err != nil {
		return models.User{}, err
	}
	newUser.ID = uint32(res.InsertedID.(int64))
	return newUser, nil
}

func (m *MongoStorage) Remove(id uint32) (uint32, error) {
	var result bson.M
	m.mu.Lock()
	defer m.mu.Unlock()
	err := m.MongoStorage.FindOneAndDelete(context.Background(), bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *MongoStorage) Update(id uint32, name, email string, age uint8) (models.User, error) {
	var result bson.M
	m.mu.Lock()
	defer m.mu.Unlock()
	user := bson.D{
		{"$set", bson.D{{"name", name}}},
		{"$set", bson.D{{"email", email}}},
		{"$set", bson.D{{"age", age}}},
	}
	err := m.MongoStorage.FindOneAndUpdate(context.Background(), bson.D{{"_id", id}}, user).Decode(&result)
	if err != nil {
		return models.User{}, err
	}
	return models.User{ID: id, Name: name, Email: email, Age: age}, nil
}
