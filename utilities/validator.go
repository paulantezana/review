package utilities

import "regexp"

func ValidateDni(dni string) bool {
	match, _ := regexp.MatchString("^[0-9]{8}$", dni)
	return match
}

func ValidateRuc() {

}
