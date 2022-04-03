package authenticationstore_test

import (
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/golib/web"
	"github.com/lonepeon/golib/web/authenticationstore"
)

func TestInMemoryAuthenticate(t *testing.T) {
	type Testcase struct {
		Username       string
		Password       string
		ExpectedID     string
		ExpectedStatus bool
	}

	inmemory := authenticationstore.NewInMemory()
	_, err := inmemory.Register("jdoe", "jdoe1234")
	testutils.RequireNoError(t, err, "unexpected error")
	_, err = inmemory.Register("jane", "jane1234")
	testutils.RequireNoError(t, err, "unexpected error")

	runner := func(name string, tc Testcase) {
		t.Run(name, func(t *testing.T) {
			id, status := inmemory.Authenticate(tc.Username, tc.Password)
			testutils.AssertEqualString(t, tc.ExpectedID, id, "unexpected user ID")
			testutils.AssertEqualBool(t, tc.ExpectedStatus, status, "unexpected authentication status")
		})
	}

	runner("when credentials are valid", Testcase{
		Username: "jdoe", Password: "jdoe1234",
		ExpectedID: "1", ExpectedStatus: true,
	})

	runner("when username doesn't match password", Testcase{
		Username: "jdoe", Password: "jane1234",
		ExpectedID: "", ExpectedStatus: false,
	})
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
	testutils.AssertEqualString(t, "jdoe1234", user.Password, "unexpected password")
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
