package response

// HealthResponse is the health check response.
type HealthResponse struct {
	Status  string         `json:"status"`
	Services []ServiceHealth `json:"services,omitempty"`
}

// ServiceHealth represents the health of a single service.
type ServiceHealth struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
}
