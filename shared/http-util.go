package shared

import "net/http"

// Response api response structure
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// setResponseHeaders set response headers like: Content-Type, Status code
func setResponseHeaders(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}
