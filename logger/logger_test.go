package logger_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/lonepeon/golib/logger"
	"github.com/lonepeon/golib/testutils"
)

type OutputMock struct {
	strings.Builder
}

func (o *OutputMock) Assert(t *testing.T, f func(*testing.T, []map[string]interface{})) {
	var logs []map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(o.String()))
	for decoder.More() {
		var log map[string]interface{}
		err := decoder.Decode(&log)
		testutils.AssertNoError(t, err, "can't decode log line")

		logs = append(logs, log)
	}

	f(t, logs)

}

func TestLoggerInInfoLevel(t *testing.T) {
	var output OutputMock
	log, closer := logger.NewLogger(&output)

	log.WithFields(logger.String("key", "value"), logger.Int("other", 42)).Info("an info message")
	log.Infof("a message with %d variable", 1)

	testutils.RequireNoError(t, closer(), "can't close logger")
	output.Assert(t, func(t *testing.T, logs []map[string]interface{}) {
		testutils.AssertEqualInt(t, 2, len(logs), "unexpected number of logs")

		testutils.AssertEqualString(t, "info", logs[0]["level"].(string), "unexpected level field: %v", logs[0])
		testutils.AssertEqualString(t, "an info message", logs[0]["msg"].(string), "unexpected msg: %v", logs[0])
		testutils.AssertEqualString(t, "value", logs[0]["key"].(string), "unexpected key field: %v", logs[0])
		testutils.AssertEqualFloat64(t, 42, logs[0]["other"].(float64), "unexpected key field: %v", logs[0])

		testutils.AssertEqualString(t, "info", logs[1]["level"].(string), "unexpected level field: %v", logs[1])
		testutils.AssertEqualString(t, "a message with 1 variable", logs[1]["msg"].(string), "unexpected msg: %v", logs[1])
	})
}

func TestLoggerInErrorLevel(t *testing.T) {
	var output OutputMock
	log, closer := logger.NewLogger(&output)

	log.
		WithFields(logger.String("some", "value"), logger.Int("else", 1337)).
		Error("an error message")
	log.Errorf("%d error message", 1)

	testutils.RequireNoError(t, closer(), "can't close logger")
	output.Assert(t, func(t *testing.T, logs []map[string]interface{}) {
		testutils.AssertEqualInt(t, 2, len(logs), "unexpected number of logs")

		testutils.AssertEqualString(t, "error", logs[0]["level"].(string), "unexpected level field: %v", logs[0])
		testutils.AssertEqualString(t, "an error message", logs[0]["msg"].(string), "unexpected msg: %v", logs[0])
		testutils.AssertEqualString(t, "value", logs[0]["some"].(string), "unexpected key field: %v", logs[0])
		testutils.AssertEqualFloat64(t, 1337, logs[0]["else"].(float64), "unexpected key field: %v", logs[0])

		testutils.AssertEqualString(t, "error", logs[1]["level"].(string), "unexpected level field: %v", logs[1])
		testutils.AssertEqualString(t, "1 error message", logs[1]["msg"].(string), "unexpected msg: %v", logs[1])
	})
}

func TestLoggerWithField(t *testing.T) {
	var output OutputMock
	log, closer := logger.NewLogger(&output)

	log.Info("default logger")
	log2 := log.WithFields(logger.String("my-field", "my-field-value"))
	log2.WithFields(logger.Int("other", 42)).Info("an info message")
	log2.Info("another info message")
	log.WithFields(logger.Int("other", 1337)).Info("another default logging with no field")

	testutils.RequireNoError(t, closer(), "can't close logger")
	output.Assert(t, func(t *testing.T, logs []map[string]interface{}) {
		testutils.AssertEqualInt(t, 4, len(logs), "unexpected number of logs")

		testutils.AssertEqualNil(t, logs[0]["my-field"], "unexpected my-field in %v", logs[0])
		testutils.AssertNotEqualNil(t, logs[1]["my-field"], "unexpected my-field in %v", logs[1])
		testutils.AssertNotEqualNil(t, logs[2]["my-field"], "unexpected my-field in %v", logs[2])
		testutils.AssertEqualNil(t, logs[3]["my-field"], "unexpected my-field in %v", logs[3])
	})
}
