package user

import (
	"context"
	"errors"
	"net/http"
	"sync"
)

type service struct {
	users *dataStore
}

// CreateUser add a new user in storage
//
// Check if user already exists, if exists then return error
func (s *service) CreateUser(ctx context.Context, req *http.Request, userReq User) (User, error) {
	s.users.mtx.Lock()
	defer s.users.mtx.Unlock()

	if _, ok := s.users.store[userReq.ID]; ok {
		return User{}, errors.New("user already exists")
	}

	s.users.store[userReq.ID] = userReq
	return s.users.store[userReq.ID], nil
}

// GetUsers get users from storage
func (s *service) GetUsers(ctx context.Context, req *http.Request) []User {
	s.users.mtx.RLock()
	defer s.users.mtx.RUnlock()

	var (
		users = make([]User, 0, len(s.users.store))
	)

	for _, v := range s.users.store {
		users = append(users, v)
	}

	return users
}

// GetUser from storage
func (s *service) GetUser(ctx context.Context, req *http.Request, userId int) (User, error) {
	s.users.mtx.RLock()
	defer s.users.mtx.RUnlock()

	if v, ok := s.users.store[userId]; ok {
		return v, nil
	}

	return User{}, errors.New("user not found")
}

// UpdateUser update user by user id
func (s *service) UpdateUser(ctx context.Context, req *http.Request, userReq User, userId int) (User, error) {
	s.users.mtx.Lock()
	defer s.users.mtx.Unlock()

	if _, ok := s.users.store[userId]; ok {
		s.users.store[userId] = User{
			ID:          userId,
			Name:        userReq.Name,
			Email:       userReq.Email,
			PhoneNumber: userReq.PhoneNumber,
		}

		return s.users.store[userId], nil
	}

	return User{}, errors.New("user not found")
}

// DeleteUser delete method is idempotent, no need to check for user existence
func (s *service) DeleteUser(ctx context.Context, req *http.Request, userId int) string {
	s.users.mtx.Lock()
	defer s.users.mtx.Unlock()

	delete(s.users.store, userId)
	return "user deleted successfully"
}

// Service on user struct
type Service interface {
	CreateUser(ctx context.Context, req *http.Request, userReq User) (User, error)
	GetUsers(ctx context.Context, req *http.Request) []User
	GetUser(ctx context.Context, req *http.Request, userId int) (User, error)
	UpdateUser(ctx context.Context, req *http.Request, userReq User, userId int) (User, error)
	DeleteUser(ctx context.Context, req *http.Request, userId int) string
}

// NewService initialises the user service
//
// Implements all the methods
func NewService() Service {
	return &service{
		users: &dataStore{
			store: make(dataMap),
			mtx:   new(sync.RWMutex),
		},
	}
}
