package storage

import (
	"database/sql"
	"fmt"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/models"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"sync"
)

type PostgreStorage struct {
	db *sql.DB
	mu sync.Mutex
}

func NewPostgreStorage() *PostgreStorage {
	user := os.Getenv("POSTGRE_USER")
	pass := os.Getenv("POSTGRE_PASS")
	host := os.Getenv("POSTGRE_HOST")
	port, err := strconv.Atoi(os.Getenv("POSTGRE_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	db := os.Getenv("POSTGRE_DB")
	postgreURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, db)

	conn, err := sql.Open("postgres", postgreURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}
	return &PostgreStorage{db: conn}
}

func (p *PostgreStorage) GetAll() ([]models.User, error) {
	rows, err := p.db.Query(`SELECT id, name, email, age FROM "user"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (p *PostgreStorage) Retrieve(id uint32) (models.User, error) {
	var user models.User
	err := p.db.QueryRow(
		`SELECT id, name, email, age FROM "user" WHERE id=$1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Age)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (p *PostgreStorage) Add(name, email string, age uint8) (models.User, error) {
	var id uint32
	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.db.QueryRow(
		`INSERT INTO "user"(name, email, age) VALUES ($1, $2, $3) RETURNING id`,
		name,
		email,
		age,
	).Scan(&id)
	if err != nil {
		return models.User{}, err
	}

	return models.User{ID: id, Name: name, Email: email, Age: age}, nil
}

func (p *PostgreStorage) Remove(id uint32) (uint32, error) {
	var _id uint32
	p.mu.Lock()
	defer p.mu.Unlock()
	err := p.db.QueryRow(
		`DELETE FROM "user" WHERE id=$1 RETURNING id`,
		id,
	).Scan(&_id)
	if err != nil {
		return 0, err
	}

	return _id, nil
}

func (p *PostgreStorage) Update(id uint32, name, email string, age uint8) (models.User, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, err := p.db.Exec(
		`UPDATE "user" SET name=$1, email=$2, age=$3 WHERE id=$4 RETURNING id`,
		name,
		email,
		age,
		id,
	)
	if err != nil {
		return models.User{}, err
	}

	return models.User{ID: id, Name: name, Email: email, Age: age}, nil
}
