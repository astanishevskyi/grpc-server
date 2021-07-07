package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/models"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type ElasticStorage struct {
	client *elastic.Client
	mu     sync.Mutex
	lastID uint32
}

func NewElasticStorage() *ElasticStorage {
	user := os.Getenv("ELASTIC_USER")
	pass := os.Getenv("ELASTIC_PASS")
	addr := os.Getenv("ELASTIC_ADDR")

	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{addr},
		Username:  user,
		Password:  pass,
	})
	if err != nil {
		log.Fatal(err)
	}
	if _, err := client.Ping(); err != nil {
		log.Fatal(err)
	}
	e := &ElasticStorage{client: client}
	lastID, _ := e.getLastID()
	e.lastID = lastID
	return e
}

func (e *ElasticStorage) getLastID() (uint32, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"aggs": map[string]interface{}{
			"max_id": map[string]interface{}{
				"max": map[string]interface{}{
					"field": "id",
				},
			},
		},
		"size": 0,
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return 0, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(context.Background()),
		e.client.Search.WithIndex("user"),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithTrackTotalHits(true),
		e.client.Search.WithPretty(),
	)
	if err != nil {
		return 0, err
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return 0, err
	}
	id := r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)

	return uint32(id), nil
}

func (e *ElasticStorage) GetAll() ([]models.User, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	var buf bytes.Buffer
	res, err := e.client.Search(
		e.client.Search.WithContext(context.Background()),
		e.client.Search.WithIndex("user"),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithTrackTotalHits(true),
		e.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	var users []models.User
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		id, err := strconv.Atoi(hit.(map[string]interface{})["_id"].(string))
		if err != nil {
			return nil, err
		}
		data := hit.(map[string]interface{})["_source"]
		name := data.(map[string]interface{})["name"].(string)
		email := data.(map[string]interface{})["email"].(string)
		age := data.(map[string]interface{})["age"].(float64)

		user := models.User{ID: uint32(id), Name: name, Email: email, Age: uint8(age)}
		users = append(users, user)
	}

	return users, nil
}

func (e *ElasticStorage) Retrieve(id uint32) (models.User, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"_id": strconv.Itoa(int(id)),
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return models.User{}, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(context.Background()),
		e.client.Search.WithIndex("user"),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithTrackTotalHits(true),
		e.client.Search.WithPretty(),
	)
	if err != nil {
		return models.User{}, err
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return models.User{}, err
	}

	var user models.User
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		id, err := strconv.Atoi(hit.(map[string]interface{})["_id"].(string))
		if err != nil {
			return models.User{}, err
		}
		data := hit.(map[string]interface{})["_source"]
		name := data.(map[string]interface{})["name"].(string)
		email := data.(map[string]interface{})["email"].(string)
		age := data.(map[string]interface{})["age"].(float64)

		user = models.User{ID: uint32(id), Name: name, Email: email, Age: uint8(age)}
	}

	return user, nil
}

func (e *ElasticStorage) Add(name, email string, age uint8) (models.User, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.lastID++
	user := models.User{Name: name, Email: email, Age: age}
	dataJSON, err := json.Marshal(user)
	if err != nil {
		return models.User{}, nil
	}

	request := esapi.IndexRequest{Index: "user", DocumentID: strconv.Itoa(int(e.lastID)), Body: strings.NewReader(string(dataJSON))}
	_, err = request.Do(context.Background(), e.client)
	if err != nil {
		return models.User{}, err
	}
	return models.User{ID: e.lastID, Name: name, Email: email, Age: age}, nil
}

func (e *ElasticStorage) Remove(id uint32) (uint32, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.client.Delete("user", strconv.Itoa(int(id)))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (e *ElasticStorage) Update(id uint32, name, email string, age uint8) (models.User, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	inputData := map[string]interface{}{
		"doc": map[string]interface{}{
			"name":  name,
			"email": email,
			"age":   age,
		},
	}
	dataJSON, err := json.Marshal(inputData)
	if err != nil {
		return models.User{}, nil
	}
	_, err = e.client.Update("user", strconv.Itoa(int(id)), strings.NewReader(string(dataJSON)))
	if err != nil {
		return models.User{}, err
	}
	return models.User{ID: id, Name: name, Email: email, Age: age}, nil
}
