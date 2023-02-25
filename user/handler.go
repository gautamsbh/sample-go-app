package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gautamsbh/sample-go-app/shared"
)

type handler struct {
	service Service
}

// Handler user request handler
type Handler interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// createUser create user handler parse the request
//
// calls user service to execute business logic
func (h handler) createUser(ctx context.Context, req *http.Request) shared.Response {
	var (
		user User
	)

	// decode json request body
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		return shared.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// call CreateUser service and read the error
	out, err := h.service.CreateUser(ctx, req, user)
	if err != nil {
		return shared.Response{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}

	return shared.Response{
		Code: http.StatusOK,
		Data: out,
	}
}

// getUsers return list of users
func (h handler) getUsers(ctx context.Context, req *http.Request) shared.Response {
	out := h.service.GetUsers(ctx, req)
	return shared.Response{
		Code: http.StatusOK,
		Data: out,
	}
}

// getUser return user detail
func (h handler) getUser(ctx context.Context, req *http.Request, userID int) shared.Response {
	out, err := h.service.GetUser(ctx, req, userID)
	if err != nil {
		return shared.Response{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}

	return shared.Response{
		Code: http.StatusOK,
		Data: out,
	}
}

// updateUser update user detail
func (h handler) updateUser(ctx context.Context, req *http.Request, userID int) shared.Response {
	var (
		user User
	)

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		return shared.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	out, err := h.service.UpdateUser(ctx, req, user, userID)
	if err != nil {
		return shared.Response{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}

	return shared.Response{
		Code: http.StatusOK,
		Data: out,
	}
}

// deleteUser delete a user from storage
func (h handler) deleteUser(ctx context.Context, req *http.Request, userID int) shared.Response {
	out := h.service.DeleteUser(ctx, req, userID)
	return shared.Response{
		Code: http.StatusOK,
		Data: out,
	}
}

// Route method acts as a router for users endpoint
func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		ctx   = req.Context()
		paths = strings.Split(strings.Trim(req.URL.Path, "/"), "/")
		out   shared.Response
	)

	if len(paths) == 1 {
		switch req.Method {
		case http.MethodGet:
			out = h.getUsers(ctx, req)
		case http.MethodPost:
			out = h.createUser(ctx, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	} else if len(paths) == 2 {
		userID, err := strconv.Atoi(paths[1])
		if err != nil {
			http.Error(w, "user id path variable invalid", http.StatusBadRequest)
			return
		}

		switch req.Method {
		case http.MethodGet:
			out = h.getUser(ctx, req, userID)
		case http.MethodPut:
			out = h.updateUser(ctx, req, userID)
		case http.MethodDelete:
			out = h.deleteUser(ctx, req, userID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	} else {
		http.NotFound(w, req)
		return
	}

	// marshal the json response
	respBody, err := json.Marshal(out)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// write response header, code and body: header should be set before code
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(out.Code)
	_, _ = w.Write(respBody)
}

// NewHandler initialises the user handler
//
// Implements all the methods
func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}
