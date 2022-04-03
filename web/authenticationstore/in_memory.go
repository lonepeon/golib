package authenticationstore

import (
	"strconv"

	"github.com/lonepeon/golib/web"
)

type InMemory struct {
	lastid int
	users  []web.User
}

func NewInMemory() *InMemory {
	return &InMemory{}
}

func (i *InMemory) Authenticate(username string, password string) (string, bool) {
	for _, user := range i.users {
		if user.Username == username && user.Password == password {
			return user.ID, true
		}
	}

	return "", false
}

func (i *InMemory) Lookup(id string) (web.User, error) {
	for _, user := range i.users {
		if user.ID == id {
			return user, nil
		}
	}

	return web.User{}, web.ErrUserNotFound
}

func (i *InMemory) Register(username string, password string) (string, error) {
	for _, user := range i.users {
		if user.Username == username {
			return "", web.ErrUserAlreadyExist
		}
	}

	i.lastid += 1

	user := web.User{ID: strconv.Itoa(i.lastid), Username: username, Password: password}

	i.users = append(i.users, user)

	return user.ID, nil
}
