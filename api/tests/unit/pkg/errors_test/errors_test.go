package errors_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	err := apperrors.New(http.StatusBadRequest, "Bad Request")
	assert.Equal(t, "Bad Request", err.Error())

	errWithDetail := apperrors.NewWithDetail(http.StatusNotFound, "Not Found", "User 42 not found")
	assert.Equal(t, "Not Found: User 42 not found", errWithDetail.Error())
}

func TestAppError_Fields(t *testing.T) {
	err := apperrors.New(http.StatusBadRequest, "Bad Request")
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	assert.Equal(t, "Bad Request", err.Title)
	assert.Equal(t, "", err.Detail)
}

func TestAppError_WithDetail(t *testing.T) {
	err := apperrors.NewWithDetail(http.StatusNotFound, "Not Found", "User 42 not found")
	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.Equal(t, "Not Found", err.Title)
	assert.Equal(t, "User 42 not found", err.Detail)
}

func TestValidationError(t *testing.T) {
	fieldErrors := map[string][]string{
		"name": {"The name field is required."},
		"type": {"The selected type is invalid."},
	}

	err := apperrors.NewValidationError(fieldErrors)
	assert.Equal(t, "The given data was invalid.", err.Error())
	assert.Equal(t, fieldErrors, err.FieldErrors)

	data, jsonErr := json.Marshal(err)
	assert.NoError(t, jsonErr)
	assert.Contains(t, string(data), `"message":"The given data was invalid."`)
	assert.Contains(t, string(data), `"name":["The name field is required."]`)
}

func TestPredefinedErrors(t *testing.T) {
	assert.Equal(t, http.StatusBadRequest, apperrors.ErrBadRequest.StatusCode)
	assert.Equal(t, http.StatusUnauthorized, apperrors.ErrUnauthorized.StatusCode)
	assert.Equal(t, http.StatusForbidden, apperrors.ErrForbidden.StatusCode)
	assert.Equal(t, http.StatusNotFound, apperrors.ErrNotFound.StatusCode)
	assert.Equal(t, http.StatusMethodNotAllowed, apperrors.ErrMethodNotAllowed.StatusCode)
	assert.Equal(t, http.StatusNotAcceptable, apperrors.ErrNotAcceptable.StatusCode)
	assert.Equal(t, http.StatusConflict, apperrors.ErrConflict.StatusCode)
	assert.Equal(t, http.StatusGone, apperrors.ErrGone.StatusCode)
	assert.Equal(t, http.StatusUnsupportedMediaType, apperrors.ErrUnsupportedMedia.StatusCode)
	assert.Equal(t, http.StatusInternalServerError, apperrors.ErrInternal.StatusCode)
	assert.Equal(t, http.StatusTooManyRequests, apperrors.ErrTooManyRequests.StatusCode)
}

func TestNotFoundResource(t *testing.T) {
	err := apperrors.NotFoundResource("wallet", 42)
	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.Equal(t, "Resource not found", err.Title)
	assert.Contains(t, err.Detail, "wallet")
	assert.Contains(t, err.Detail, "42")
	assert.Equal(t, "NotFoundHttpException", err.Exception)
}

func TestGoneTransaction(t *testing.T) {
	err := apperrors.GoneTransaction()
	assert.Equal(t, http.StatusGone, err.StatusCode)
	assert.Contains(t, err.Detail, "200032")
	assert.Contains(t, err.Detail, "rule deleted")
}

func TestDemoUserBlocked(t *testing.T) {
	err := apperrors.DemoUserBlocked()
	assert.Equal(t, http.StatusForbidden, err.StatusCode)
	assert.Contains(t, err.Detail, "Demo user")
}

func TestNotAcceptableHeader(t *testing.T) {
	err := apperrors.NotAcceptableHeader("text/xml")
	assert.Equal(t, http.StatusNotAcceptable, err.StatusCode)
	assert.Contains(t, err.Detail, "text/xml")
}

func TestUnsupportedContentType(t *testing.T) {
	err := apperrors.UnsupportedContentType("text/xml")
	assert.Equal(t, http.StatusUnsupportedMediaType, err.StatusCode)
	assert.Contains(t, err.Detail, "text/xml")
}

func TestEmptyContentType(t *testing.T) {
	err := apperrors.EmptyContentType()
	assert.Equal(t, http.StatusUnsupportedMediaType, err.StatusCode)
	assert.Contains(t, err.Detail, "cannot be empty")
}

func TestErrors_As(t *testing.T) {
	var appErr *apperrors.AppError
	var valErr *apperrors.ValidationError

	// Test AppError As
	err1 := apperrors.New(http.StatusBadRequest, "test")
	assert.True(t, errors.As(err1, &appErr))
	assert.False(t, errors.As(err1, &valErr))

	// Test ValidationError As
	err2 := apperrors.NewValidationError(map[string][]string{"field": {"error"}})
	assert.True(t, errors.As(err2, &valErr))
	assert.False(t, errors.As(err2, &appErr))
}

func TestAppError_JSONSerialization(t *testing.T) {
	err := apperrors.NewWithDetail(http.StatusNotFound, "Resource not found", "User 42 not found")
	data, jsonErr := json.Marshal(err)
	assert.NoError(t, jsonErr)
	assert.Contains(t, string(data), `"status":404`)
	assert.Contains(t, string(data), `"title":"Resource not found"`)
	assert.Contains(t, string(data), `"detail":"User 42 not found"`)
}
