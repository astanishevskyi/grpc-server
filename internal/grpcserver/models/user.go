package models

type User struct {
	ID    uint32 `json:"id"`
	Age   uint8  `json:"age"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserService interface {
	GetAll() []User
	Retrieve(id uint32) (User, error)
	Add(name, email string, age uint8) (User, error)
	Remove(id uint32) (uint32, error)
	Update(id uint32, name, email string, age uint8) (User, error)
}
