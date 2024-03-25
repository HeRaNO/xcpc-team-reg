package utils

import (
	"regexp"
	"strings"
)

var numberReg = regexp.MustCompile(`^\d+$`)

func IsNumber(x *string) bool {
	return numberReg.Match([]byte(*x))
}

func TrimName(x *string) *string {
	r := strings.TrimSpace(*x)
	return &r
}
