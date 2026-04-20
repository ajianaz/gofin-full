package domain

import "time"

type Rule struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	UserGroupID int64    `json:"user_group_id" db:"user_group_id"`
	RuleGroupID *int64   `json:"rule_group_id,omitempty" db:"rule_group_id"`
	Title      string    `json:"title" db:"title"`
	Description *string  `json:"description,omitempty" db:"description"`
	Priority   int       `json:"priority" db:"priority"`
	Active     bool      `json:"active" db:"active"`
	Strict     bool      `json:"strict" db:"strict"`
	StopProcessing bool  `json:"stop_processing" db:"stop_processing"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"-" db:"deleted_at"`

	// Joined
	Triggers []RuleTrigger `json:"triggers,omitempty" db:"-"`
	Actions  []RuleAction  `json:"actions,omitempty" db:"-"`
	Group    *RuleGroup    `json:"group,omitempty" db:"-"`
}

type RuleGroup struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	UserGroupID int64    `json:"user_group_id" db:"user_group_id"`
	Title      string    `json:"title" db:"title"`
	Active     bool      `json:"active" db:"active"`
	Order      int       `json:"order" db:"order"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"-" db:"deleted_at"`
}

type RuleTrigger struct {
	ID         int64     `json:"id" db:"id"`
	RuleID     int64     `json:"rule_id" db:"rule_id"`
	TriggerType string   `json:"trigger_type" db:"trigger_type"`
	TriggerValue string  `json:"trigger_value" db:"trigger_value"`
	StopProcessing bool `json:"stop_processing" db:"stop_processing"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type RuleAction struct {
	ID        int64     `json:"id" db:"id"`
	RuleID    int64     `json:"rule_id" db:"rule_id"`
	ActionType string   `json:"action_type" db:"action_type"`
	ActionValue string  `json:"action_value" db:"action_value"`
	Order     int       `json:"order" db:"order"`
	StopProcessing bool `json:"stop_processing" db:"stop_processing"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
