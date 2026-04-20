package request

// LoginRequest for email + password authentication.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest for self-registration (when AUTH_ALLOW_REGISTRATION=true).
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// AdminCreateUserRequest for admin user creation.
type AdminCreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// OAuthCallbackRequest for OAuth code exchange.
type OAuthCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state"`
}

// RefreshTokenRequest for token refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UpdateUserRequest for profile updates.
type UpdateUserRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
	Name  string `json:"name" validate:"omitempty,min=1"`
}

// ChangePasswordRequest for password changes.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=16"`
}

// CreateUserGroupRequest for group creation.
type CreateUserGroupRequest struct {
	Title string `json:"title" validate:"required,min=1"`
}

// UpdateUserGroupRequest for group updates.
type UpdateUserGroupRequest struct {
	Title string `json:"title" validate:"required,min=1"`
}

// SwitchGroupRequest for group switching.
type SwitchGroupRequest struct {
	UserGroupID int64 `json:"user_group_id" validate:"required"`
}

// CreateWalletRequest for wallet creation.
type CreateWalletRequest struct {
	Name            string   `json:"name" validate:"required,min=1"`
	AccountType     string   `json:"account_type" validate:"required"`
	IBAN            string   `json:"iban,omitempty"`
	BIC             string   `json:"bic,omitempty"`
	CurrencyID      string   `json:"currency_id,omitempty"`
	Active          *bool    `json:"active,omitempty"`
	IncludeNetWorth *bool    `json:"include_net_worth,omitempty"`
	Latitude        *float64 `json:"latitude,omitempty"`
	Longitude       *float64 `json:"longitude,omitempty"`
	Notes           string   `json:"notes,omitempty"`
}

// UpdateWalletRequest for wallet updates.
type UpdateWalletRequest struct {
	Name            string   `json:"name,omitempty" validate:"omitempty,min=1"`
	IBAN            string   `json:"iban,omitempty"`
	BIC             string   `json:"bic,omitempty"`
	CurrencyID      string   `json:"currency_id,omitempty"`
	Active          *bool    `json:"active,omitempty"`
	IncludeNetWorth *bool    `json:"include_net_worth,omitempty"`
	Latitude        *float64 `json:"latitude,omitempty"`
	Longitude       *float64 `json:"longitude,omitempty"`
	Notes           string   `json:"notes,omitempty"`
}
