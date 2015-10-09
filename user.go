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

func FindUser(username, password string) (*User, error) {
	out := &User{
		Username: username,
	}

	existingUser, err := globalUserStore.FindByUsername(username)
	switch {
	case err != nil:
		return out, err
	case existingUser == nil:
		return out, errCredentialsIncorrect
	case bcrypt.CompareHashAndPassword(
		[]byte(existingUser.HashedPassword),
		[]byte(password),
	) != nil:
		return out, errCredentialsIncorrect
	}

	return existingUser, nil
}

func UpdateUser(user *User, email, currentPassword, newPassword string) (User, error) {
	out := *user
	out.Email = email

	// Check if the email exists
	existingUser, err := globalUserStore.FindByEmail(email)
	switch {
	case err != nil:
		return out, err
	case existingUser != nil && existingUser.ID != user.ID:
		return out, errEmailExists
	}

	// At this point, we can update the email address
	user.Email = email

	// No current password? Don't try to update the password
	switch {
	case currentPassword == "":
		return out, nil
	case bcrypt.CompareHashAndPassword(
		[]byte(user.HashedPassword),
		[]byte(currentPassword),
	) != nil:
		return out, errPasswordIncorrect
	case newPassword == "":
		return out, errNoPassword
	case len(newPassword) < passwordLength:
		return out, errPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), hashCost)
	user.HashedPassword = string(hashedPassword)
	return out, err
}
