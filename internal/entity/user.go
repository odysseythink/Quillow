package entity

import "time"

type User struct {
	ID            uint
	UserGroupID   uint
	Email         string
	Password      string
	RememberToken string
	Reset         string
	Blocked       bool
	BlockedCode   string
	Domain        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Role struct {
	ID          uint
	Name        string
	DisplayName string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GroupMembership struct {
	ID          uint
	UserID      uint
	UserGroupID uint
	UserRoleID  uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type InvitedUser struct {
	ID         uint
	UserID     uint
	Email      string
	InviteCode string
	Expires    time.Time
	Redeemed   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
