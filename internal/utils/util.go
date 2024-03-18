package utils

import (
	"html/template"
	"regexp"
	"strings"
)

var numberReg = regexp.MustCompile(`^\d+$`)

func IsNumber(x *string) bool {
	return numberReg.Match([]byte(*x))
}

func TrimName(x *string) *string {
	r := strings.TrimSpace(*x)
	r = template.HTMLEscapeString(r)
	return &r
}
