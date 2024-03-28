package utils

import "regexp"

var numberReg = regexp.MustCompile(`^\d+$`)

func IsNumber(x *string) bool {
	return numberReg.Match([]byte(*x))
}
