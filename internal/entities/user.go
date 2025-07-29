package entities

import (
	"net/mail"
	"scheduler/internal/errors"
	"scheduler/pkg/utils"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string
	Username  string
	Email     string
	password  string
	Credits   int
	CreatedAt string
	UpdatedAt string
	hashed    bool
}

func NewUser(username string, email string, password string) (*User, error) {
	if !validPassword(password) {
		return nil, errors.PASSWORD_INVALID()
	}

	if !validUsername(username) {
		return nil, errors.USERNAME_INVALID()
	}

	if !validEmail(email) {
		return nil, errors.EMAIL_INVALID()
	}

	now := time.Now().Format(time.RFC3339)
	hashed, err := hashPassword(password)

	if err != nil {
		return nil, errors.PASSWORD_HASHING()
	}

	return &User{
		Id:        uuid.NewString(),
		Username:  username,
		Email:     email,
		password:  hashed,
		Credits:   0,
		CreatedAt: now,
		UpdatedAt: now,
		hashed:    true,
	}, nil
}

func (u *User) GetPassword() (string, error) {
	if u.hashed {
		return u.password, nil
	}

	hashed, err := hashPassword(u.password)

	if err != nil {
		return "", errors.PASSWORD_HASHING()
	}

	u.password = hashed
	u.hashed = true

	return u.password, nil
}

func (u *User) CheckPasswordHash(password string) bool {
	return checkPasswordHash(password, u.password)
}

func validUsername(username string) bool {
	cleanedStr := strings.Join(strings.Fields(username), " ")

	if len(cleanedStr) < 5 {
		return false
	}

	if len(cleanedStr) > 20 {
		return false
	}

	for _, c := range cleanedStr {
		if len(strings.TrimSpace(string(c))) == 0 {
			continue
		}

		if !utils.IsAlphanumeric(byte(c)) {
			return false
		}
	}

	return true
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func validPassword(password string) bool {
	letters := 0
	upper := false
	lower := false
	symbol := false
	number := false

	for _, c := range password {
		if unicode.IsUpper(c) {
			upper = true
		} else if unicode.IsLower(c) {
			lower = true
		} else if unicode.IsNumber(c) {
			number = true
		} else {
			symbol = true
		}

		letters += 1
	}

	return letters >= 8 && number && upper && lower && symbol
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
