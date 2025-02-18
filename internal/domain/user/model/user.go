package model

import (
	"errors"
	"regexp"
	"time"
)

type Subscribtion uint

const (
	Normal Subscribtion = iota
	Premium1
	Premium2
)

type (
	UserID uint
	Phone  string
)

func (p Phone) IsValid() bool {
	re := regexp.MustCompile(`^\+989\d{9}$`)
	return re.MatchString(string(p))
}

func (u *User) Validate() error {
	if !u.Phone.IsValid() {
		return errors.New("invalid phone")
	}

	if u.FirstName == "" {
		return errors.ErrUnsupported
	}

	if u.LastName == "" {
		return errors.ErrUnsupported
	}

	return nil
}

type User struct {
	UserID           UserID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	FirstName        string
	LastName         string
	Phone            Phone
	SubscribtionType Subscribtion
}

type TokenBlacklist struct {
	ID        uint
	CreatedAt time.Time
	ExpiresAt time.Time
	Value     string
	UserID    UserID
}
