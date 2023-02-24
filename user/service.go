package user

import (
	"context"
	"errors"
	"net/http"
)

type service struct {
	users UsersStorage
}

// CreateUser add a new user in storage
//
// Check if user already exists, if exists then return error
func (s service) CreateUser(ctx context.Context, req *http.Request, userReq User) (User, error) {
	if _, ok := s.users[userReq.ID]; ok {
		return User{}, errors.New("user already exists")
	}

	s.users[userReq.ID] = userReq
	return s.users[userReq.ID], nil
}

// GetUsers get users from storage
func (s service) GetUsers(ctx context.Context, req *http.Request) []User {
	var (
		users = make([]User, 0, len(s.users))
	)

	for _, v := range s.users {
		users = append(users, v)
	}

	return users
}

// GetUser from storage
func (s service) GetUser(ctx context.Context, req *http.Request, userId int) (User, error) {
	if v, ok := s.users[userId]; ok {
		return v, nil
	}

	return User{}, errors.New("user not found")
}

// UpdateUser update user by user id
func (s service) UpdateUser(ctx context.Context, req *http.Request, userReq User, userId int) (User, error) {
	if _, ok := s.users[userId]; ok {
		s.users[userId] = User{
			ID:          userId,
			Name:        userReq.Name,
			Email:       userReq.Email,
			PhoneNumber: userReq.PhoneNumber,
		}

		return s.users[userId], nil
	}

	return User{}, errors.New("user not found")
}

// DeleteUser delete method is idempotent, no need to check for user existence
func (s service) DeleteUser(ctx context.Context, req *http.Request, userId int) string {
	delete(s.users, userId)
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
		users: make(UsersStorage),
	}
}
