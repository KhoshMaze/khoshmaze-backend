package model

import (
	"errors"
	"regexp"
	"strings"
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
	ID               UserID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	FirstName        string
	LastName         string
	Phone            Phone
	SubscribtionType Subscribtion
}

type TokenBlacklist struct {
	ExpiresAt time.Time
	Value     string
	UserID    UserID
}

type UserFilter struct {
	ID    UserID
	Phone string
}

func (f *UserFilter) IsValid() bool {
	f.Phone = strings.TrimSpace(f.Phone)
	return f.ID > 0 || len(f.Phone) > 0
}
