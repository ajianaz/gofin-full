package response

// Envelope is the standard API response wrapper.
type Envelope struct {
	Data  interface{} `json:"data,omitempty"`
	Meta  *Meta       `json:"meta,omitempty"`
	Links *Links      `json:"links,omitempty"`
}

// ErrorEnvelope is the standard error response.
type ErrorEnvelope struct {
	Message   string            `json:"message,omitempty"`
	Exception string            `json:"exception,omitempty"`
	Errors    map[string][]string `json:"errors,omitempty"`
}

// Meta contains pagination and other metadata.
type Meta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination contains pagination info.
type Pagination struct {
	Total       int64  `json:"total"`
	Count       int     `json:"count"`
	PerPage     int     `json:"per_page"`
	CurrentPage int     `json:"current_page"`
	TotalPages  int     `json:"total_pages"`
}

// Links contains navigation links.
type Links struct {
	Self    string `json:"self,omitempty"`
	Next    string `json:"next,omitempty"`
	Prev    string `json:"prev,omitempty"`
}

// NewEnvelope creates a response envelope with data.
func NewEnvelope(data interface{}) *Envelope {
	return &Envelope{Data: data}
}

// NewPaginatedEnvelope creates a response envelope with pagination.
func NewPaginatedEnvelope(data interface{}, total int64, page, perPage int) *Envelope {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &Envelope{
		Data: data,
		Meta: &Meta{
			Pagination: &Pagination{
				Total:       total,
				Count:       perPage,
				PerPage:     perPage,
				CurrentPage: page,
				TotalPages:  totalPages,
			},
		},
	}
}

// NewErrorEnvelope creates an error response.
func NewErrorEnvelope(message, exception string) *ErrorEnvelope {
	return &ErrorEnvelope{
		Message:   message,
		Exception: exception,
	}
}

// NewValidationErrorEnvelope creates a validation error response.
func NewValidationErrorEnvelope(errors map[string][]string) *ErrorEnvelope {
	return &ErrorEnvelope{
		Message: "The given data was invalid.",
		Errors:  errors,
	}
}
