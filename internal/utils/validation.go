package utils

import (
	"errors"
	"regexp"
	"unicode"
)

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasMinLen || !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("password must contain at least 8 characters, upper, lower, number, and special symbol")
	}

	return nil
}

func ValidateName(name string) error {
	re := regexp.MustCompile(`^[A-Za-zА-Яа-яЁё\s\-]+$`)
	if len(name) == 0 {
		return errors.New("name cannot be empty")
	}
	if !re.MatchString(name) {
		return errors.New("name contains invalid characters")
	}
	return nil
}
