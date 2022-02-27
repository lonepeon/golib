package env_test

import (
	"testing"

	"github.com/lonepeon/golib/env"
	"github.com/lonepeon/golib/testutils"
)

type Cfg struct {
	KeyString         string   `env:"KEY_STRING"`
	KeyStringRequired string   `env:"KEY_STRING_REQUIRED,required=true"`
	KeyStringDefault  string   `env:"KEY_STRING_DEFAULT,default=default value"`
	KeyInt            int      `env:"KEY_INT"`
	KeyIntRequired    int      `env:"KEY_INT_REQUIRED,required=true"`
	KeyIntDefault     int      `env:"KEY_INT_DEFAULT,default=1337"`
	KeyStrings        []string `env:"KEY_STRINGS"`
	KeyStrings2       []string `env:"KEY_STRINGS2,sep=|"`
}

func TestLoadSuccess(t *testing.T) {
	var cfg Cfg

	t.Setenv("KEY_STRING", "value")
	t.Setenv("KEY_STRING_DEFAULT", "")
	t.Setenv("KEY_STRING_REQUIRED", "required value")
	t.Setenv("KEY_INT", "42")
	t.Setenv("KEY_INT_REQUIRED", "-5")
	t.Setenv("KEY_INT_DEFAULT", "")
	t.Setenv("KEY_STRINGS", "aaaa,bb,cccc")
	t.Setenv("KEY_STRINGS2", "dddd|ee|ffff")

	err := env.Load(&cfg)
	testutils.AssertNoError(t, err, "expected to load config")
	testutils.AssertEqualString(t, "value", cfg.KeyString, "expected to load cfg.KeyString")
	testutils.AssertEqualString(t, "required value", cfg.KeyStringRequired, "expected to load cfg.KeyStringRequired")
	testutils.AssertEqualString(t, "default value", cfg.KeyStringDefault, "expected to load cfg.KeyStringDefault")
	testutils.AssertEqualInt(t, 42, cfg.KeyInt, "expected to load cfg.KeyInt")
	testutils.AssertEqualInt(t, -5, cfg.KeyIntRequired, "expected to load cfg.KeyIntRequired")
	testutils.AssertEqualInt(t, 1337, cfg.KeyIntDefault, "expected to load cfg.KeyIntDefault")
	testutils.AssertEqualStrings(t, []string{"aaaa", "bb", "cccc"}, cfg.KeyStrings, "expected to load cfg.KeyStrings")
	testutils.AssertEqualStrings(t, []string{"dddd", "ee", "ffff"}, cfg.KeyStrings2, "expected to load cfg.KeyStrings2")
}
