package user

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/gautamsbh/sample-go-app/shared"
)

var (
	usersWithIdRoute = regexp.MustCompile(`^\/users\/(\d+)$`)
)

type handler struct {
	service Service
}

// Handler user request handler
type Handler interface {
	RegisterRoutes(gRouter shared.GenericRouter)
}

// createUser create user handler parse the request
//
// calls user service to execute business logic
func (h *handler) createUser(ctx context.Context, req *http.Request) shared.Response {
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
func (h *handler) getUsers(ctx context.Context, req *http.Request) shared.Response {
	out := h.service.GetUsers(ctx, req)
	return shared.Response{
		Code: http.StatusOK,
		Data: out,
	}
}

// getUser return user detail
func (h *handler) getUser(ctx context.Context, req *http.Request) shared.Response {
	userID, err := shared.ParseUserIDFromPath(usersWithIdRoute, req.URL.Path)
	if err != nil {
		return shared.Response{
			Code:    http.StatusUnprocessableEntity,
			Message: "user id param invalid",
		}
	}

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
func (h *handler) updateUser(ctx context.Context, req *http.Request) shared.Response {
	var (
		user User
	)

	userID, err := shared.ParseUserIDFromPath(usersWithIdRoute, req.URL.Path)
	if err != nil {
		return shared.Response{
			Code:    http.StatusUnprocessableEntity,
			Message: "user id param invalid",
		}
	}

	err = json.NewDecoder(req.Body).Decode(&user)
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
func (h *handler) deleteUser(ctx context.Context, req *http.Request) shared.Response {
	userID, err := shared.ParseUserIDFromPath(usersWithIdRoute, req.URL.Path)
	if err != nil {
		return shared.Response{
			Code:    http.StatusUnprocessableEntity,
			Message: "user id param invalid",
		}
	}

	out := h.service.DeleteUser(ctx, req, userID)
	return shared.Response{
		Code: http.StatusOK,
		Data: out,
	}
}

// RegisterRoutes register all user routes
func (h *handler) RegisterRoutes(gRouter shared.GenericRouter) {
	gRouter.Post(`^\/users[\/]*$`, h.createUser)
	gRouter.Get(`^\/users[\/]*$`, h.getUsers)
	gRouter.Get(`^\/users\/(\d+)$`, h.getUser)
	gRouter.Put(`^\/users\/(\d+)$`, h.updateUser)
	gRouter.Delete(`^\/users\/(\d+)$`, h.deleteUser)
}

// NewHandler initialises the user handler
//
// Implements all the methods
func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}
