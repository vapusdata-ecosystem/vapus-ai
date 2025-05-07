package prompts

import (
	"fmt"
	"regexp"
	"strings"
)

func BuildTagPattern(tag string) (string, string, *regexp.Regexp) {
	regex := regexp.MustCompile(fmt.Sprintf(Baseregex, tag, tag))
	return strings.Replace(StartTagTemplate, "TAG", tag, -1),
		strings.Replace(EndTagTemplate, "TAG", tag, -1), regex

}
