package validator

import "testing"

func TestValidators(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		ok   bool
	}{
		{"profile ok", ValidateProfileName("work-1"), true},
		{"profile path", ValidateProfileName("../evil"), false},
		{"email ok", ValidateEmail("a.b+c@example.co"), true},
		{"email bad", ValidateEmail("a@b"), false},
		{"full name newline", ValidateFullName("a\nb"), false},
		{"full name ok", ValidateFullName("Ana María"), true},
		{"username space", ValidateUsername("a b"), false},
	} {
		if (tc.err == nil) != tc.ok {
			t.Errorf("%s: err = %v", tc.name, tc.err)
		}
	}
}
