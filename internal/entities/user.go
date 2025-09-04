package entities

import (
	"net/mail"
	"scheduler/internal/errors"
	"scheduler/pkg/utils"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseEntity
	username       string
	email          string
	password       string
	credits        int
	frozen_credits int
}

func NewUser(username string, email string, password string) (*User, error) {
	user := &User{
		BaseEntity:     *NewBaseEntity(),
		credits:        0,
		frozen_credits: 0,
	}

	err := user.SetUsername(username)
	if err != nil {
		return nil, err
	}

	err = user.SetEmail(email)
	if err != nil {
		return nil, err
	}

	err = user.SetPassword(password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func HydrateUser(id string, createdAt, updatedAt time.Time,
	username, email, hashedPassword string,
	credits, frozenCredits int,
) *User {
	return &User{
		BaseEntity: BaseEntity{
			id:        id,
			createdAt: createdAt,
			updatedAt: updatedAt,
		},
		username:       username,
		email:          email,
		password:       hashedPassword,
		credits:        credits,
		frozen_credits: frozenCredits,
	}
}

func (u *User) SetUsername(username string) error {
	cleanedStr := strings.Join(strings.Fields(username), " ")

	if len(cleanedStr) < 5 || len(cleanedStr) > 20 {
		return errors.INVALID_FIELD_VALUE("username")
	}

	for _, c := range cleanedStr {
		if len(strings.TrimSpace(string(c))) == 0 {
			continue
		}

		if !utils.IsAlphanumeric(byte(c)) {
			return errors.INVALID_FIELD_VALUE("username")
		}
	}

	u.username = cleanedStr

	return nil
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) SetEmail(email string) error {
	mail, err := mail.ParseAddress(email)
	if err != nil {
		return errors.INVALID_FIELD_VALUE("email")
	}

	u.email = mail.Address

	return nil
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) SetPassword(password string) error {
	count := 0
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

		count += 1
	}

	if count < 8 || !number || !upper || !lower || !symbol {
		return errors.INVALID_FIELD_VALUE("password")
	}

	hash, err := hashPassword(password)
	if err != nil {
		return errors.PASSWORD_HASHING()
	}

	u.password = hash

	return nil
}

func (u *User) GetPassword() (string, error) {
	_, err := bcrypt.Cost([]byte(u.password))
	if err == nil {
		return u.password, err
	}

	hashed, err := hashPassword(u.password)
	if err != nil {
		return "", errors.PASSWORD_HASHING()
	}

	u.password = hashed

	return u.password, nil
}

func (u *User) AddCredits(amount int) {
	u.credits += amount
}

func (u *User) RemoveCredits(amount int) error {
	total := u.credits - amount

	if total < 0 {
		return errors.NOT_ENOUGH_CREDITS_ERROR()
	}

	u.credits = total

	return nil
}

func (u *User) GetCredits() int {
	return u.credits
}

func (u *User) AddFrozenCredits(amount int) error {
	err := u.RemoveCredits(amount)
	if err != nil {
		return err
	}

	u.frozen_credits += amount

	return nil
}

func (u *User) RemoveFrozenCredits(amount int, refund bool) error {
	if u.frozen_credits-amount < 0 {
		return errors.NOT_ENOUGH_CREDITS_ERROR()
	}

	u.frozen_credits -= amount

	if refund {
		u.credits += amount
	}

	return nil
}

func (u *User) GetFrozenCredits() int {
	return u.frozen_credits
}

func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))

	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
