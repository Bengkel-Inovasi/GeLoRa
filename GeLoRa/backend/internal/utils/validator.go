package utils

import (
	"errors"
	"regexp"
	"unicode"

	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func ValidateName(name string) error {
	if len(name) < 3 || len(name) > 128 {
		return errors.New("Name must be between 3 and 128 characters")
	}
	return nil
}

func ValidateUsername(username string) error {
	if len(username) < 8 || len(username) > 64 {
		return errors.New("Username must be between 8 and 64 characters")
	}
	if !usernameRegex.MatchString(username) {
		return errors.New("Username may only contain alphanumerics, dashes, and underscores")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters")
	}
	var hasLetter, hasNumber bool
	for _, c := range password {
		switch {
		case unicode.IsLetter(c):
			hasLetter = true
		case unicode.IsNumber(c):
			hasNumber = true
		}
	}
	if !hasLetter || !hasNumber {
		return errors.New("Password must contain at least one letter and one number")
	}
	return nil
}

func ValidateBio(bio string) error {
	if len(bio) > 2048 {
		return errors.New("Bio must not exceed 2048 characters")
	}
	return nil
}

func ValidateRole(role string) error {
	switch domainmodel.UserRole(role) {
	case domainmodel.UserRoleSuper, domainmodel.UserRoleAdmin, domainmodel.UserRoleClient:
		return nil
	}
	return errors.New("Role must be one of: super, admin, client")
}
