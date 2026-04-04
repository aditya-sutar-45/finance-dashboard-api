// Package models
package models

import (
	"time"

	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginUserResponse struct {
	SessionID             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  User      `json:"user"`
}

type RenewAcessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func DatabaseUserToUser(u database.User) User {
	var updatedAt time.Time
	var createdAt time.Time

	if u.UpdatedAt.Valid {
		updatedAt = u.UpdatedAt.Time
	}

	if u.CreatedAt.Valid {
		createdAt = u.CreatedAt.Time
	}

	user := User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return user
}

func DatabaseUsersToUsers(dbUsers []database.User) []User {
	var users []User
	for _, u := range dbUsers {
		user := DatabaseUserToUser(u)
		users = append(users, user)
	}

	return users
}
