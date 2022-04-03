package authenticationstore_test

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3" // sqlite3 adapter

	"github.com/lonepeon/golib/sqlutil"
	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/golib/web"
	"github.com/lonepeon/golib/web/authenticationstore"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Parallel()

	t.Run("SQLiteAuthenticateSuccess", testSQLiteAuthenticateSuccess)
	t.Run("SQLiteAuthenticateErrorWrongPassword", testSQLiteAuthenticateErrorWrongPassword)
	t.Run("SQLiteAuthenticateErrorWrongUsername", testSQLiteAuthenticateErrorWrongUsername)
	t.Run("SQLiteLookupSuccess", testSQLiteLookupSuccess)
	t.Run("SQLiteLookupError", testSQLiteLookupError)
}

func testSQLiteLookupSuccess(t *testing.T) {
	store := setupDatabase(t)
	id, err := store.Register("jane", "jane1234")
	testutils.RequireNoError(t, err, "can't insert user in database")

	user, err := store.Lookup(id)
	testutils.RequireNoError(t, err, "didn't find user")
	testutils.AssertEqualString(t, id, user.ID, "unexpected ID")
	testutils.AssertEqualString(t, "jane", user.Username, "unexpected username")
}

func testSQLiteLookupError(t *testing.T) {
	store := setupDatabase(t)
	_, err := store.Register("jane", "jane1234")
	testutils.RequireNoError(t, err, "can't insert user in database")

	user, err := store.Lookup("04663061-16c3-425f-84d1-96cf027f275f")
	testutils.RequireHasError(t, err, "didn't expect to find user: %v", user)
	testutils.AssertErrorIs(t, web.ErrUserNotFound, err, "unexpected error")
}

func testSQLiteAuthenticateSuccess(t *testing.T) {
	store := setupDatabase(t)
	_, err := store.Register("jane", "jane1234")
	testutils.RequireNoError(t, err, "can't insert user in database")

	id, err := store.Authenticate("jane", "jane1234")
	testutils.RequireNoError(t, err, "unexpected authentication error")
	testutils.AssertNotEmptyString(t, id, "expecting an id")
}

func testSQLiteAuthenticateErrorWrongPassword(t *testing.T) {
	store := setupDatabase(t)
	_, err := store.Register("jane", "jane1234")
	testutils.RequireNoError(t, err, "can't insert user in database")

	id, err := store.Authenticate("jane", "jane")
	testutils.RequireHasError(t, err, "expecting authentication error but got an ID: %v", id)
	testutils.AssertErrorIs(t, web.ErrUserInvalidCredentials, err, "unexpected error")
	testutils.AssertContainsString(t, "invalid password", err.Error(), "unexpected error message")
}

func testSQLiteAuthenticateErrorWrongUsername(t *testing.T) {
	store := setupDatabase(t)
	_, err := store.Register("jane", "jane1234")
	testutils.RequireNoError(t, err, "can't insert user in database")

	id, err := store.Authenticate("jdoe", "jane1234")
	testutils.RequireHasError(t, err, "expecting authentication error but got an ID: %v", id)
	testutils.AssertErrorIs(t, web.ErrUserInvalidCredentials, err, "unexpected error")
	testutils.AssertContainsString(t, "invalid username", err.Error(), "unexpected error message")
}

func setupDatabase(t *testing.T) *authenticationstore.SQLite {
	f, err := ioutil.TempFile("", "authentication-*.sqlite")
	testutils.RequireNoError(t, err, "can't create SQLite temporary file")

	db, err := sql.Open("sqlite3", f.Name())
	testutils.RequireNoError(t, err, "can't open sqlite connection")

	_, err = sqlutil.ExecuteMigrations(context.Background(), db, authenticationstore.Migrations())
	testutils.RequireNoError(t, err, "can't run migrations")

	t.Cleanup(func() {
		db.Close()
		os.Remove(f.Name())
	})

	return authenticationstore.NewSQLite(db, "a-pepper")
}
