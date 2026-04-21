package domain

import (
	"github.com/google/uuid"
	"time"
)

type Webhook struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	UserGroupID uuid.UUID `json:"user_group_id" db:"user_group_id"`
	Title       string    `json:"title" db:"title"`
	URL         string    `json:"url" db:"url"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`

	// Joined
	Triggers []WebhookTrigger `json:"triggers,omitempty" db:"-"`
}

type WebhookTrigger struct {
	ID       uuid.UUID `json:"id" db:"id"`
	WebhookID uuid.UUID `json:"webhook_id" db:"webhook_id"`
	Trigger   string `json:"trigger" db:"trigger"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type WebhookMessage struct {
	ID        uuid.UUID `json:"id" db:"id"`
	WebhookID uuid.UUID `json:"webhook_id" db:"webhook_id"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type WebhookDelivery struct {
	ID              uuid.UUID `json:"id" db:"id"`
	WebhookMessageID uuid.UUID `json:"webhook_message_id" db:"webhook_message_id"`
	ResponseCode    int       `json:"response_code" db:"response_code"`
	ResponseBody    *string   `json:"response_body,omitempty" db:"response_body"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type WebhookAttempt struct {
	ID              uuid.UUID `json:"id" db:"id"`
	WebhookMessageID uuid.UUID `json:"webhook_message_id" db:"webhook_message_id"`
	Attempt         int       `json:"attempt" db:"attempt"`
	ResponseCode    int       `json:"response_code" db:"response_code"`
	ResponseBody    *string   `json:"response_body,omitempty" db:"response_body"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type WebhookResponse struct {
	ID              uuid.UUID `json:"id" db:"id"`
	WebhookMessageID uuid.UUID `json:"webhook_message_id" db:"webhook_message_id"`
	ResponseCode    int       `json:"response_code" db:"response_code"`
	ResponseBody    *string   `json:"response_body,omitempty" db:"response_body"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
