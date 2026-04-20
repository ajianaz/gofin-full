package domain

import "time"

type Attachment struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	AttachableType string `json:"attachable_type" db:"attachable_type"`
	AttachableID   int64  `json:"attachable_id" db:"attachable_id"`
	Filename   string    `json:"filename" db:"filename"`
	MimeType   string    `json:"mime_type" db:"mime_type"`
	Size       int64     `json:"size" db:"size"`
	Uploaded   bool      `json:"uploaded" db:"uploaded"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"-" db:"deleted_at"`
}
