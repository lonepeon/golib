package cli

import "strings"

type ArgStrings []string

func (a ArgStrings) String() string {
	return strings.Join(a, ", ")
}

func (a *ArgStrings) Set(s string) error {
	*a = append(*a, s)

	return nil
}
