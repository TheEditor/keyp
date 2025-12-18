package store

import "errors"

var (
	ErrNotFound        = errors.New("secret not found")
	ErrAlreadyExists   = errors.New("secret already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrDatabaseLocked  = errors.New("database is locked")
	ErrVaultClosed     = errors.New("vault is closed")
)
