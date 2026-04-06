package entity

import "time"

type RuleGroup struct {
	ID             uint
	UserID         uint
	UserGroupID    uint
	Title          string
	Description    string
	Order          uint
	Active         bool
	StopProcessing bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type Rule struct {
	ID             uint
	UserID         uint
	UserGroupID    uint
	RuleGroupID    uint
	Title          string
	Description    string
	Order          uint
	Active         bool
	StopProcessing bool
	Strict         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type RuleTrigger struct {
	ID             uint
	RuleID         uint
	TriggerType    string
	TriggerValue   string
	Order          uint
	Active         bool
	StopProcessing bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RuleAction struct {
	ID             uint
	RuleID         uint
	ActionType     string
	ActionValue    string
	Order          uint
	Active         bool
	StopProcessing bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
