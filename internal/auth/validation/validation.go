package validation

import (
	"regexp"
)

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^.+@.+$`)
	return re.MatchString(email)
}

func ValidateLogin(login string) bool {
	re := regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-_.]{3,16}[A-Za-z0-9]$`)
	return re.MatchString(login)
}

func ValidatePassword(password string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+=-]{8,16}$`)
	return re.MatchString(password)
}

func ValidateName(name string) bool {
	re := regexp.MustCompile(`^.{5,50}$`)
	return re.MatchString(name)
}
