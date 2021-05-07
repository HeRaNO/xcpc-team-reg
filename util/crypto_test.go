package util_test

import (
	"testing"
	"unicode"

	"github.com/HeRaNO/xcpc-team-reg/util"
)

func TestGenToken(t *testing.T) {
	p := 20
	token, err := util.GenToken(p)
	if err != nil {
		t.Fatalf("Cannot generate token. Len: %d, err: %s", p, err.Error())
	}
	if len(token) != p {
		t.Fatalf("token %s: length is not %d", token, p)
	}
	for _, let := range token {
		if !unicode.IsLetter(let) && !unicode.IsDigit(let) {
			t.Fatalf("token %s: has other letter", token)
		}
	}
}

func BenchmarkGenToken(b *testing.B) {
	p := 20
	for i := 0; i < b.N; i++ {
		_, err := util.GenToken(p)
		if err != nil {
			b.Fatalf("Cannot generate token. Len: %d, err: %s", p, err.Error())
		}
	}
}
