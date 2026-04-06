package model

import "time"

type RuleGroupModel struct {
	ID             uint       `gorm:"primaryKey;column:id"`
	UserID         uint       `gorm:"column:user_id"`
	UserGroupID    uint       `gorm:"column:user_group_id"`
	Title          string     `gorm:"column:title"`
	Description    *string    `gorm:"column:description"`
	Order          uint       `gorm:"column:order"`
	Active         bool       `gorm:"column:active"`
	StopProcessing bool       `gorm:"column:stop_processing"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
}

func (RuleGroupModel) TableName() string { return "rule_groups" }

type RuleModel struct {
	ID             uint       `gorm:"primaryKey;column:id"`
	UserID         uint       `gorm:"column:user_id"`
	UserGroupID    uint       `gorm:"column:user_group_id"`
	RuleGroupID    uint       `gorm:"column:rule_group_id"`
	Title          string     `gorm:"column:title"`
	Description    *string    `gorm:"column:description"`
	Order          uint       `gorm:"column:order"`
	Active         bool       `gorm:"column:active"`
	StopProcessing bool       `gorm:"column:stop_processing"`
	Strict         bool       `gorm:"column:strict"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
}

func (RuleModel) TableName() string { return "rules" }

type RuleTriggerModel struct {
	ID             uint      `gorm:"primaryKey;column:id"`
	RuleID         uint      `gorm:"column:rule_id"`
	TriggerType    string    `gorm:"column:trigger_type"`
	TriggerValue   string    `gorm:"column:trigger_value"`
	Order          uint      `gorm:"column:order"`
	Active         bool      `gorm:"column:active"`
	StopProcessing bool      `gorm:"column:stop_processing"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (RuleTriggerModel) TableName() string { return "rule_triggers" }

type RuleActionModel struct {
	ID             uint      `gorm:"primaryKey;column:id"`
	RuleID         uint      `gorm:"column:rule_id"`
	ActionType     string    `gorm:"column:action_type"`
	ActionValue    string    `gorm:"column:action_value"`
	Order          uint      `gorm:"column:order"`
	Active         bool      `gorm:"column:active"`
	StopProcessing bool      `gorm:"column:stop_processing"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (RuleActionModel) TableName() string { return "rule_actions" }
