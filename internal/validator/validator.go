package validator

import (
    "fmt"
    "regexp"
    "strings"
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
    emailRegex      = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
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
    
    return nil
}

// ValidateSigningKey checks if the signing key is valid (optional)
func ValidateSigningKey(key string) error {
    key = strings.TrimSpace(key)
    
    // Signing key is optional, so empty is valid
    if len(key) == 0 {
        return nil
    }
    
    // Basic GPG key format validation (simplified)
    if len(key) < 16 {
        return fmt.Errorf("signing key must be at least 16 characters long")
    }
    
    return nil
}