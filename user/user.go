package user

import "sync"

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type (
	dataMap map[int]User
)

// DataStore is to store user information and support concurrency
type dataStore struct {
	store dataMap
	mtx   *sync.RWMutex
}
