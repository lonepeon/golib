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

	inmemory := authenticationstore.NewInMemory([]web.User{
		{Username: "jdoe", Password: "jdoe1234"},
		{Username: "jane", Password: "jane1234"},
	})

	runner := func(name string, tc Testcase) {
		t.Run(name, func(t *testing.T) {
			id, status := inmemory.Authenticate(tc.Username, tc.Password)
			testutils.AssertEqualString(t, tc.ExpectedID, id, "unexpected user ID")
			testutils.AssertEqualBool(t, tc.ExpectedStatus, status, "unexpected authentication status")
		})
	}

	runner("when credentials are valid", Testcase{
		Username: "jdoe", Password: "jdoe1234",
		ExpectedID: "jdoe", ExpectedStatus: true,
	})

	runner("when username doesn't match password", Testcase{
		Username: "jdoe", Password: "jane1234",
		ExpectedID: "", ExpectedStatus: false,
	})
}

func TestInMemoryLookupSuccess(t *testing.T) {
	inmemory := authenticationstore.NewInMemory([]web.User{
		{Username: "jdoe", Password: "jdoe1234"},
		{Username: "jane", Password: "jane1234"},
	})

	user, err := inmemory.Lookup("jdoe")
	testutils.RequireNoError(t, err, "unexpected error")
	testutils.AssertEqualString(t, "jdoe", user.Username, "unexpected user name")
	testutils.AssertEqualString(t, "jdoe1234", user.Password, "unexpected password")
}

func TestInMemoryLookupError(t *testing.T) {
	inmemory := authenticationstore.NewInMemory([]web.User{
		{Username: "jdoe", Password: "jdoe1234"},
		{Username: "jane", Password: "jane1234"},
	})

	_, err := inmemory.Lookup("msmith")
	testutils.AssertErrorIs(t, err, web.ErrUserNotFound, "unexpected error")
}
