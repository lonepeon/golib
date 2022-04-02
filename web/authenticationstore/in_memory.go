package authenticationstore

import (
	"github.com/lonepeon/golib/web"
)

type InMemory struct {
	users []web.User
}

func NewInMemory(users []web.User) InMemory {
	return InMemory{users: users}
}

func (i InMemory) Authenticate(username string, password string) (string, bool) {
	for _, user := range i.users {
		if user.Username == username && user.Password == password {
			return user.Username, true
		}
	}

	return "", false
}

func (i InMemory) Lookup(id string) (web.User, error) {
	for _, user := range i.users {
		if user.Username == id {
			return user, nil
		}
	}

	return web.User{}, web.ErrUserNotFound
}
