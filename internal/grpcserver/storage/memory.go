package storage

import (
	"errors"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/models"
	"sync"
)

type UserStorage struct {
	mu         sync.Mutex
	lastID     uint32
	UsersSlice []models.User `json:" users"`
}

func NewInMemoryUserStorage() *UserStorage {
	storage := UserStorage{lastID: 1, UsersSlice: make([]models.User, 0)}
	return &storage
}

func (s *UserStorage) GetAll() ([]models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	copyUserSlice := make([]models.User, len(s.UsersSlice))
	copy(copyUserSlice, s.UsersSlice)
	return copyUserSlice, nil
}

func (s *UserStorage) Retrieve(id uint32) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, val := range s.UsersSlice {
		if val.ID == id {
			return val, nil
		}
	}
	return models.User{}, errors.New("no user was found")
}

func (s *UserStorage) Add(name, email string, age uint8) (models.User, error) {
	user := models.User{ID: s.lastID, Age: age, Name: name, Email: email}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastID++
	s.UsersSlice = append(s.UsersSlice, user)
	return user, nil
}

func (s *UserStorage) Remove(id uint32) (uint32, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, val := range s.UsersSlice {
		if val.ID == id {
			s.UsersSlice = append(s.UsersSlice[:i], s.UsersSlice[i+1:]...)
			return id, nil
		}
	}
	return 0, errors.New("no user was found")
}

func (s *UserStorage) Update(id uint32, name, email string, age uint8) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, val := range s.UsersSlice {
		if val.ID == id {
			val.Name = name
			val.Email = email
			val.Age = age
			s.UsersSlice[i] = val
			return val, nil
		}
	}
	return models.User{}, errors.New("no user was found")
}
