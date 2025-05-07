package dmutils

import "regexp"

var TemplateVarRegex = regexp.MustCompile(`\{\{(.*?)\}\}`)
