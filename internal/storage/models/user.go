package models

import (
	"database/sql"
	"time"
)

type Role struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique"`
	Users []User `gorm:"foreignKey:RoleID"`
}

type User struct {
	ID          int64  `gorm:"primaryKey"`
	Username    string `gorm:"unique"`
	CreatedAt   time.Time
	Description sql.NullString
	RoleID      uint
	Role        Role
}
