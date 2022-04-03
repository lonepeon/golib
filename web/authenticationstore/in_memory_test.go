package authenticationstore_test

import (
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/golib/web"
	"github.com/lonepeon/golib/web/authenticationstore"
)

func TestInMemoryAuthenticateSuccess(t *testing.T) {
	inmemory := authenticationstore.NewInMemory()
	_, err := inmemory.Register("jdoe", "jdoe1234")
	testutils.RequireNoError(t, err, "unexpected error")
	_, err = inmemory.Register("jane", "jane1234")
	testutils.RequireNoError(t, err, "unexpected error")

	id, err := inmemory.Authenticate("jdoe", "jdoe1234")
	testutils.RequireNoError(t, err, "unexpected error")
	testutils.AssertEqualString(t, "1", id, "unexpected user ID")
}

func TestInMemoryAuthenticateError(t *testing.T) {
	inmemory := authenticationstore.NewInMemory()
	_, err := inmemory.Register("jdoe", "jdoe1234")
	testutils.RequireNoError(t, err, "unexpected error")
	_, err = inmemory.Register("jane", "jane1234")
	testutils.RequireNoError(t, err, "unexpected error")

	id, err := inmemory.Authenticate("jdoe", "jane1234")
	testutils.RequireHasError(t, err, "expecting an error but got an ID: %v", id)
	testutils.AssertErrorIs(t, web.ErrUserInvalidCredentials, err, "unexpected error")
}

func TestInMemoryLookupSuccess(t *testing.T) {
	inmemory := authenticationstore.NewInMemory()
	id, err := inmemory.Register("jdoe", "jdoe1234")
	testutils.RequireNoError(t, err, "unexpected error")
	_, err = inmemory.Register("jane", "jane1234")
	testutils.RequireNoError(t, err, "unexpected error")

	user, err := inmemory.Lookup(id)
	testutils.RequireNoError(t, err, "unexpected error")
	testutils.AssertEqualString(t, "1", user.ID, "unexpected user ID")
	testutils.AssertEqualString(t, "jdoe", user.Username, "unexpected user name")
}

func TestInMemoryLookupError(t *testing.T) {
	inmemory := authenticationstore.NewInMemory()
	_, err := inmemory.Register("jdoe", "jdoe1234")
	testutils.RequireNoError(t, err, "unexpected error")
	_, err = inmemory.Register("jane", "jane1234")
	testutils.RequireNoError(t, err, "unexpected error")

	_, err = inmemory.Lookup("msmith")
	testutils.AssertErrorIs(t, err, web.ErrUserNotFound, "unexpected error")
}

func TestInMemoryRegisterSuccess(t *testing.T) {
	inmemory := authenticationstore.NewInMemory()

	id, err := inmemory.Register("jdoe", "jdoe1234")
	testutils.RequireNoError(t, err, "unexpected error")
	testutils.AssertEqualString(t, "1", id, "unexpected identifier")

	id, err = inmemory.Register("jane", "jane1234")
	testutils.RequireNoError(t, err, "unexpected error")
	testutils.AssertEqualString(t, "2", id, "unexpected identifier")
}

func TestInMemoryRegisterError(t *testing.T) {
	inmemory := authenticationstore.NewInMemory()

	id, err := inmemory.Register("jdoe", "jdoe1234")
	testutils.RequireNoError(t, err, "unexpected error")
	testutils.AssertEqualString(t, "1", id, "unexpected identifier")

	_, err = inmemory.Register("jdoe", "jane1234")
	testutils.RequireHasError(t, err, "expecting an error")
	testutils.AssertErrorIs(t, web.ErrUserAlreadyExist, err, "wrong error")
}
