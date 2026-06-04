package response

// HealthResponse is the health check response.
type HealthResponse struct {
	Status        string          `json:"status"`
	Version       string          `json:"version,omitempty"`
	UptimeSeconds int64           `json:"uptime_seconds,omitempty"`
	Services      []ServiceHealth `json:"services,omitempty"`
}

// ServiceHealth represents the health of a single service.
type ServiceHealth struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
}
