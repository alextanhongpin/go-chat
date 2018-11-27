package entity

import (
	"errors"
	"time"

	"github.com/alextanhongpin/go-openid/pkg/passwd"
)

var ErrUserNotFound = errors.New("user not found")

// User represents the user of the application.
type User struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at,omitempty"`
	HashedPassword string    `json:"hashed_password"`
}

// NewUser returns a new user with the given name and email.
func NewUser(name, email string) *User {
	return &User{
		Name:  name,
		Email: email,
	}
}

// SetPassword takes a plaintext password and sets the hashed password.
func (u *User) SetPassword(password string) error {
	var err error
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	u.HashedPassword, err = passwd.Hash(password)
	return err
}

// ComparePassword checks if the password matches the given user.
func (u *User) ComparePassword(password string) error {
	return passwd.Verify(password, u.HashedPassword)
}

// func (u *User) Save(repo repository.User) error {
//         var err error
//         u.ID, err = repo.Create(u)
//         return err
// }
