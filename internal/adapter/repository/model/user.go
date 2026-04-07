package model

import "time"

type UserModel struct {
	ID            uint      `gorm:"primaryKey;column:id"`
	UserGroupID   *uint     `gorm:"column:user_group_id"`
	Email         string    `gorm:"column:email"`
	Password      string    `gorm:"column:password"`
	RememberToken *string   `gorm:"column:remember_token"`
	Reset         *string   `gorm:"column:reset"`
	Blocked       bool      `gorm:"column:blocked"`
	BlockedCode   *string   `gorm:"column:blocked_code"`
	Domain        *string   `gorm:"column:domain"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (UserModel) TableName() string { return "users" }

type RoleModel struct {
	ID          uint      `gorm:"primaryKey;column:id"`
	Name        string    `gorm:"column:name"`
	DisplayName *string   `gorm:"column:display_name"`
	Description *string   `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (RoleModel) TableName() string { return "roles" }

type RoleUserModel struct {
	UserID uint `gorm:"primaryKey;column:user_id"`
	RoleID uint `gorm:"primaryKey;column:role_id"`
}

func (RoleUserModel) TableName() string { return "role_user" }

type GroupMembershipModel struct {
	ID          uint       `gorm:"primaryKey;column:id"`
	UserID      uint       `gorm:"column:user_id"`
	UserGroupID uint       `gorm:"column:user_group_id"`
	UserRoleID  uint       `gorm:"column:user_role_id"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (GroupMembershipModel) TableName() string { return "group_memberships" }

type InvitedUserModel struct {
	ID         uint      `gorm:"primaryKey;column:id"`
	UserID     uint      `gorm:"column:user_id"`
	Email      string    `gorm:"column:email"`
	InviteCode string    `gorm:"column:invite_code"`
	Expires    time.Time `gorm:"column:expires"`
	Redeemed   bool      `gorm:"column:redeemed"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (InvitedUserModel) TableName() string { return "invited_users" }

type PersonalAccessTokenModel struct {
	ID            uint       `gorm:"primaryKey;column:id"`
	TokenableID   uint       `gorm:"column:tokenable_id"`
	TokenableType string     `gorm:"column:tokenable_type"`
	Name          string     `gorm:"column:name"`
	Token         string     `gorm:"column:token;type:varchar(64);uniqueIndex"`
	Abilities     *string    `gorm:"column:abilities"`
	LastUsedAt    *time.Time `gorm:"column:last_used_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
}

func (PersonalAccessTokenModel) TableName() string { return "personal_access_tokens" }
