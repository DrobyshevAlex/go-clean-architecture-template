package entities

import (
	"database/sql"
	"time"
)

type User struct {
	Id           int
	FirstName    string
	LastName     string
	UserName     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime //TODO: replace to custom type (wrong for clean arhitecture)
}

type UserProfile struct {
	Id        int
	FirstName string
	LastName  string
	UserName  string
}
