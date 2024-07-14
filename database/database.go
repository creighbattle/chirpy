package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
	Emails map[string]int `json:"emails"`
}

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID int `json:"id"`
	Email string `json:"email"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) CreateUser(body string, password string) (UserResponse, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return UserResponse{}, err
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return UserResponse{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:   id,
		Email: body,
		Password: string(bcryptPassword),
	}

	_, ok := dbStructure.Emails[body]
	if ok {
		return UserResponse{}, errors.New("email already exists")
	}
	dbStructure.Emails[body] = id
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{Email: body, ID: id}, nil
}

func (db *DB) UpdateUser(updatedEmail string, password string, id int) (UserResponse, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return UserResponse{}, err
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return UserResponse{}, err
	}

	allEmails := dbStructure.Emails
	user := dbStructure.Users[id]

	_, ok := allEmails[updatedEmail]
	if ok {
		return UserResponse{}, errors.New("email already exists")
	}

	delete(allEmails, user.Email)
	allEmails[updatedEmail] = user.ID

	user.Email = updatedEmail
	user.Password = string(bcryptPassword)

	dbStructure.Users[id] = user
	dbStructure.Emails = allEmails

	err = db.writeDB(dbStructure)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{Email: updatedEmail, ID: id}, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
		Emails: map[string]int{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) LoadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}