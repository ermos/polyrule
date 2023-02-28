package golang

import "regexp"

type AuthRules struct{}

func (AuthRules) GetUsernameMessage() interface{} {
	return "this username is invalid"
}

func (AuthRules) ValidateUsernameWithErrors(input string) (isValid bool, errors []string) {
	if !regexp.MustCompile("/^[a-zA-ZÀ-ÖØ-öø-ÿ ']+$/").MatchString(input) {
		errors = append(errors, "regex")
	}

	return len(errors) == 0, errors
}
func (r AuthRules) ValidateUsername(input string) (isValid bool) {
	isValid, _ = r.ValidateUsernameWithErrors(input)
	return
}
