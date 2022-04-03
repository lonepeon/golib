package authenticationstore

import (
	"strconv"

	"github.com/lonepeon/golib/web"
)

type InMemoryUser struct {
	user     web.User
	password string
}

type InMemory struct {
	lastid int
	users  []InMemoryUser
}

func NewInMemory() *InMemory {
	return &InMemory{}
}

func (i *InMemory) Authenticate(username string, password string) (string, error) {
	for _, user := range i.users {
		if user.user.Username == username && user.password == password {
			return user.user.ID, nil
		}
	}

	return "", web.ErrUserInvalidCredentials
}

func (i *InMemory) Lookup(id string) (web.User, error) {
	for _, user := range i.users {
		if user.user.ID == id {
			return user.user, nil
		}
	}

	return web.User{}, web.ErrUserNotFound
}

func (i *InMemory) Register(username string, password string) (string, error) {
	for _, user := range i.users {
		if user.user.Username == username {
			return "", web.ErrUserAlreadyExist
		}
	}

	i.lastid += 1

	user := InMemoryUser{
		user: web.User{
			ID:       strconv.Itoa(i.lastid),
			Username: username,
		},
		password: password,
	}

	i.users = append(i.users, user)

	return user.user.ID, nil
}
