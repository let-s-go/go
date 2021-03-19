package extra

import (
	"unicode"
)

func SetLowerNamingStrategy() {
	SetNamingStrategy(LowerCase)
}

func LowerCase(name string) string {
	rs := []rune(name)
	for i, r := range rs {
		if unicode.IsLower(r) {
			break
		}
		rs[i] = unicode.ToLower(r)
	}
	return string(rs)
}
