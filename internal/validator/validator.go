package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ProfileNameRules contains the validation rules for profile names
const ProfileNameRules = "Profile name must:\n" +
	"- Be between 1 and 32 characters\n" +
	"- Contain only letters, numbers, hyphens, and underscores\n" +
	"- Not start or end with hyphen or underscore"

// EmailRules contains the validation rules for email addresses
const EmailRules = "Email must:\n" +
	"- Be a valid email address\n" +
	"- Contain @ symbol\n" +
	"- Have a valid domain"

var (
	profileNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]*[a-zA-Z0-9]$|^[a-zA-Z0-9]$`)
	emailRegex       = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// ValidateProfileName checks if the profile name meets the requirements
func ValidateProfileName(name string) error {
	name = strings.TrimSpace(name)

	if len(name) == 0 {
		return fmt.Errorf("profile name cannot be empty")
	}

	if len(name) > 32 {
		return fmt.Errorf("profile name cannot be longer than 32 characters")
	}

	if !profileNameRegex.MatchString(name) {
		return fmt.Errorf("invalid profile name format. %s", ProfileNameRules)
	}

	return nil
}

// ValidateEmail checks if the email address is valid
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if len(email) == 0 {
		return fmt.Errorf("email cannot be empty")
	}

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format. %s", EmailRules)
	}

	return nil
}

// ValidateUsername checks if the username is valid
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)

	if len(username) == 0 {
		return fmt.Errorf("username cannot be empty")
	}

	if len(username) > 39 {
		return fmt.Errorf("username cannot be longer than 39 characters")
	}

	if strings.ContainsFunc(username, unicode.IsSpace) {
		return fmt.Errorf("username cannot contain whitespace")
	}

	return nil
}

// ValidateFullName checks the value used for git user.name.
func ValidateFullName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("full name cannot be empty")
	}
	if strings.ContainsFunc(name, unicode.IsControl) {
		return fmt.Errorf("full name cannot contain control characters")
	}
	return nil
}
