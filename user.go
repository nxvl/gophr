package main

import "code.google.com/p/go.crypto/bcrypt"

const (
	passwordLength = 8
	hashCost       = 10
	userIDLength   = 16
)

type User struct {
	ID             string
	Email          string
	HashedPassword string
	Username       string
}

func NewUser(username, email, password string) (User, error) {
	user := User{
		Email:    email,
		Username: username,
	}

	switch {
	case username == "":
		return user, errNoUsername
	case email == "":
		return user, errNoEmail
	case password == "":
		return user, errNoPassword
	case len(password) < passwordLength:
		return user, errPasswordTooShort
	}

	// Check if the username exists
	existingUser, err := globalUserStore.FindByUsername(username)
	switch {
	case err != nil:
		return user, err
	case existingUser != nil:
		return user, errUsernameExists
	}

	existingUser, err = globalUserStore.FindByEmail(email)
	switch {
	case err != nil:
		return user, err
	case existingUser != nil:
		return user, errEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)

	user.HashedPassword = string(hashedPassword)
	user.ID = GenerateID("usr", userIDLength)

	return user, err

}
