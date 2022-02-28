package job_test

import (
	"testing"
	"time"

	"github.com/lonepeon/golib/job"
	"github.com/lonepeon/golib/testutils"
)

func TestNextAttempt(t *testing.T) {
	j, err := job.NewJob("my-job", map[string]int{"a": 1, "b": 2})
	testutils.RequireNoError(t, err, "unexpected error while building job")
	j.MaxAttempts = 4

	testutils.AssertEqualString(t, "my-job", j.Name, "unexpected job name")

	now := time.Now()
	j2, ok := j.ConfigureNextAttempt(now)
	testutils.RequireEqualBool(t, true, ok, "unexpected job to fail at 2nd attempt")
	testutils.AssertEqualString(t, j.Name, j2.Name, "unexpected job 2 name")
	testutils.AssertEqualInt64(t, now.Add(21*time.Second).Unix(), j2.At.Unix(), "unexpected job 2 rescheduling time")

	now = time.Now()
	j3, ok := j2.ConfigureNextAttempt(now)
	testutils.RequireEqualBool(t, true, ok, "unexpected job to fail at 3rd attempt")
	testutils.AssertEqualString(t, j.Name, j3.Name, "unexpected job 3 name")
	testutils.AssertEqualInt64(t, now.Add(86*time.Second).Unix(), j3.At.Unix(), "unexpected job 3 rescheduling time")

	now = time.Now()
	j4, ok := j3.ConfigureNextAttempt(now)
	testutils.RequireEqualBool(t, true, ok, "unexpected job to fail at 4rd attempt")
	testutils.AssertEqualString(t, j.Name, j4.Name, "unexpected job 4 name")
	testutils.AssertEqualInt64(t, now.Add(261*time.Second).Unix(), j4.At.Unix(), "unexpected job 4 rescheduling time")

	now = time.Now()
	_, ok = j4.ConfigureNextAttempt(now)
	testutils.AssertEqualBool(t, false, ok, "expected job to fail after 4 attempt")
}
