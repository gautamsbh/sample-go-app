package shared

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"sync"
)

type (
	RouteHandler func(ctx context.Context, req *http.Request) Response
)

type route struct {
	method  string
	pattern *regexp.Regexp
	handler RouteHandler
}

// genericRouter generic router provide a routing layer on top of http handler
type genericRouter struct {
	routes []route
}

// GenericRouter provides functionality to granular level control over routing
//
// Implement all 5 REST methods alongwith http handler method
type GenericRouter interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	Post(path string, f RouteHandler)
	Get(path string, f RouteHandler)
	Put(path string, f RouteHandler)
	Patch(path string, f RouteHandler)
	Delete(path string, f RouteHandler)
}

// addRoute add a new route in generic router
func (g *genericRouter) addRoute(method, pattern string, f RouteHandler) {
	var (
		mtx sync.Mutex
	)

	mtx.Lock()
	defer mtx.Unlock()

	g.routes = append(g.routes, route{method: method, pattern: regexp.MustCompile(pattern), handler: f})
}

// FindRoute find route matches request
func (g *genericRouter) findRoute(req *http.Request) *route {
	for _, gRoute := range g.routes {
		if req.Method == gRoute.method && gRoute.pattern.MatchString(req.URL.Path) {
			return &gRoute
		}
	}

	return nil
}

// ServeHTTP implements server http handler
func (g *genericRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		out          Response
		ctx          = req.Context()
		matchedRoute = g.findRoute(req)
	)

	// route not found
	if matchedRoute == nil {
		http.NotFound(w, req)
		return
	}

	out = matchedRoute.handler(ctx, req)

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

// Post register post method
func (g *genericRouter) Post(pattern string, f RouteHandler) {
	g.addRoute(http.MethodPost, pattern, f)
}

// Get register get method
func (g *genericRouter) Get(pattern string, f RouteHandler) {
	g.addRoute(http.MethodGet, pattern, f)
}

// Put register put method
func (g *genericRouter) Put(pattern string, f RouteHandler) {
	g.addRoute(http.MethodPut, pattern, f)
}

// Patch register patch method
func (g *genericRouter) Patch(pattern string, f RouteHandler) {
	g.addRoute(http.MethodPatch, pattern, f)
}

// Delete register delete method
func (g *genericRouter) Delete(pattern string, f RouteHandler) {
	g.addRoute(http.MethodDelete, pattern, f)
}

// New new method initializes new generic router
func NewGenericRouter() GenericRouter {
	return new(genericRouter)
}
