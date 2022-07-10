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
		data: make(map[string]string),
	}
}

type Authors interface {
	Add(username string, email string) error
	Check(email string) bool
	Find(email string) (username *string, err error)
	Details(username string) ([]string, error)
}

type authors struct {
	data  map[string]string
	mutex sync.RWMutex
}

func (a *authors) Add(username string, email string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if username == "" || email == "" {
		return ErrAuthorCannotBeEmpty
	}

	a.data[email] = username

	return nil
}

func (a *authors) Check(email string) bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	_, ok := a.data[email]

	return ok
}

func (a *authors) Find(email string) (*string, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	v, ok := a.data[email]
	if !ok {
		return nil, ErrAuthorNotFound
	}

	return &v, nil
}

func (a *authors) Details(username string) ([]string, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	emails := make([]string, 0)
	for email, v := range a.data {
		if v == username {
			emails = append(emails, email)
		}
	}

	if len(emails) == 0 {
		return nil, ErrAuthorNotFound
	}

	return emails, nil
}
