package userDomain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string    `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Role           Role      `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPassword = string(hashedPassword)
	return nil
}

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}

func (u *User) UpdatePassword(newPassword string) error {
	err := u.VerifyPassword(newPassword)

	if err == nil {
		return ErrInvalidNewPassword
	}

	setPasswordErr := u.SetPassword(newPassword)
	if setPasswordErr != nil {
		return setPasswordErr
	}
	return nil
}
