package validators

import (
	"fmt"
	"regexp"
	"time"
	"unicode"

	"github.com/MatheusHenrique129/bemax-api/internal/core/domain"
)

var (
	phoneRegex = regexp.MustCompile(`^\d{10,14}$`)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("full name is required")
	}

	if len(name) < domain.MinNameLen {
		return fmt.Errorf("full name must be at least %d characters", domain.MinNameLen)
	}

	return nil
}

func ValidatePhone(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone is required")
	}

	phone = phoneRegex.ReplaceAllString(phone, "")

	if len(phone) < 8 || len(phone) > 15 {
		return fmt.Errorf("phone number must have between 8 and 15 digits")
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < domain.MinPasswordLen {
		return fmt.Errorf("password must be at least %d characters long", domain.MinPasswordLen)
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	isStrongPass := hasUpper && hasLower && hasDigit && hasSpecial

	if !isStrongPass {
		return fmt.Errorf("password must have at least one uppercase letter, one lowercase letter, one digit and one special character")
	}

	return nil
}

func ValidateBirthDate(birthDate string) (time.Time, error) {
	if birthDate == "" {
		return time.Time{}, fmt.Errorf("birth date is required")
	}

	birth, err := time.Parse("02/01/2006", birthDate)
	if err != nil {
		birth, err = time.Parse("2006-01-02", birthDate)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid date format, use DD/MM/YYYY or YYYY-MM-DD")
		}
	}

	if birth.After(time.Now().UTC()) {
		return time.Time{}, fmt.Errorf("birth date cannot be in the future")
	}

	maxAge := 120 * 365 * 24 * time.Hour
	if time.Since(birth) > maxAge {
		return time.Time{}, fmt.Errorf("birth date is too far in the past")
	}

	return birth, nil
}
