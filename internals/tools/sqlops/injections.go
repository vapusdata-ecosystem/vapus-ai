package sqlops

import "regexp"

func isAlphanumeric(input string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(input)
}
