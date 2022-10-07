package utils

import (
	"errors"
	"sync"
)

var _ Authors = &authors{}

var (
	ErrAuthorNotFound      = errors.New("author not found")
	ErrAuthorCannotBeEmpty = errors.New("author cannot be empty")
)

func NewAuthors() *authors {
	return &authors{
		usernames: make(map[string]string),
		upis:      make(map[string]string),
		fullnames: make(map[string]string),
	}
}

type Authors interface {
	Add(username string, email string) error
	Check(email string) bool
	FindUserName(email string) (username *string, err error)
	FindUPI(email string) (upi *string, err error)
	// full name as "first last"
	FindFullName(email string) (fullname *string, err error)
	Details(username string) ([]string, error)
}

type authors struct {
	usernames map[string]string // email => username
	upis      map[string]string // email => upi
	fullnames map[string]string // email => fullname
	mutex     sync.RWMutex
}

func (a *authors) Add(username string, email string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if username == "" || email == "" {
		return ErrAuthorCannotBeEmpty
	}

	a.usernames[email] = username

	return nil
}

func (a *authors) Check(email string) bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	_, ok := a.usernames[email]

	return ok
}

func (a *authors) FindUserName(email string) (*string, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	v, ok := a.usernames[email]
	if !ok {
		return nil, ErrAuthorNotFound
	}

	return &v, nil
}

func (a *authors) FindUPI(email string) (*string, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	v, ok := a.upis[email]
	if !ok {
		return nil, ErrAuthorNotFound
	}

	return &v, nil
}

func (a *authors) FindFullName(email string) (*string, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	v, ok := a.fullnames[email]
	if !ok {
		return nil, ErrAuthorNotFound
	}

	return &v, nil
}

func (a *authors) Details(username string) ([]string, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	emails := make([]string, 0)
	for email, v := range a.usernames {
		if v == username {
			emails = append(emails, email)
		}
	}

	if len(emails) == 0 {
		return nil, ErrAuthorNotFound
	}

	return emails, nil
}
