package response_test

import (
	"encoding/json"
	"testing"

	response "github.com/ajianaz/gofin-full/api/internal/dto/response"
	"github.com/stretchr/testify/assert"
)

func TestNewEnvelope(t *testing.T) {
	data := map[string]string{"key": "value"}
	env := response.NewEnvelope(data)

	assert.NotNil(t, env)
	assert.Equal(t, data, env.Data)
	assert.Nil(t, env.Meta)
	assert.Nil(t, env.Links)

	jsonData, err := json.Marshal(env)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"data":{"key":"value"}`)
	assert.NotContains(t, string(jsonData), `"meta"`)
}

func TestNewPaginatedEnvelope(t *testing.T) {
	items := []string{"a", "b", "c"}
	env := response.NewPaginatedEnvelope(items, 10, 1, 3)

	assert.NotNil(t, env)
	assert.NotNil(t, env.Meta)
	assert.NotNil(t, env.Meta.Pagination)

	assert.Equal(t, int64(10), env.Meta.Pagination.Total)
	assert.Equal(t, 3, env.Meta.Pagination.PerPage)
	assert.Equal(t, 1, env.Meta.Pagination.CurrentPage)
	assert.Equal(t, 4, env.Meta.Pagination.TotalPages) // ceil(10/3) = 4
}

func TestNewPaginatedEnvelope_ExactPages(t *testing.T) {
	items := []string{"a", "b"}
	env := response.NewPaginatedEnvelope(items, 6, 1, 2)

	assert.Equal(t, 3, env.Meta.Pagination.TotalPages) // 6/2 = 3 exact
}

func TestNewErrorEnvelope(t *testing.T) {
	env := response.NewErrorEnvelope("Something went wrong", "HttpException")

	jsonData, err := json.Marshal(env)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"message":"Something went wrong"`)
	assert.Contains(t, string(jsonData), `"exception":"HttpException"`)
}

func TestNewValidationErrorEnvelope(t *testing.T) {
	fieldErrors := map[string][]string{
		"name": {"The name field is required."},
	}
	env := response.NewValidationErrorEnvelope(fieldErrors)

	jsonData, err := json.Marshal(env)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"message":"The given data was invalid."`)
	assert.Contains(t, string(jsonData), `"name":["The name field is required."]`)
}

func TestHealthResponse(t *testing.T) {
	health := response.HealthResponse{
		Status: "ok",
		Services: []response.ServiceHealth{
			{Name: "postgresql", Status: "ok"},
			{Name: "redis", Status: "ok"},
		},
	}

	jsonData, err := json.Marshal(health)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"status":"ok"`)
	assert.Contains(t, string(jsonData), `"postgresql"`)
	assert.Contains(t, string(jsonData), `"redis"`)
}

func TestHealthResponse_Degraded(t *testing.T) {
	health := response.HealthResponse{
		Status: "degraded",
		Services: []response.ServiceHealth{
			{Name: "postgresql", Status: "ok"},
			{Name: "redis", Status: "error", Error: "connection refused"},
		},
	}

	jsonData, err := json.Marshal(health)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"status":"degraded"`)
	assert.Contains(t, string(jsonData), `"error":"connection refused"`)
}
