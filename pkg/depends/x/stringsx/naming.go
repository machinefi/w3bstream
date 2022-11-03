package stringsx

import (
	"strings"
	"unicode"
)

// LowerSnakeCase e.g. i_am_a_10_years_senior
func LowerSnakeCase(name string) string {
	return rewords(name, func(res, word string, idx int) string {
		if idx == 0 {
			return res + strings.ToLower(word)
		}
		return res + "_" + strings.ToLower(word)
	})
}

// UpperSnakeCase e.g. I_AM_A_10_YEARS_SENIOR
func UpperSnakeCase(name string) string {
	return rewords(name, func(res, word string, idx int) string {
		if idx == 0 {
			return res + strings.ToUpper(word)
		}
		return res + "_" + strings.ToUpper(word)
	})
}

// LowerCamelCase e.g. iAmA10YearsSenior
func LowerCamelCase(name string) string {
	return rewords(name, func(res, word string, idx int) string {
		word = strings.ToLower(word)
		runes := []rune(word)
		if idx > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		return res + string(runes)
	})
}

// UpperCamelCase e.g. IAmA10YearsSenior
func UpperCamelCase(name string) string {
	return rewords(name, func(res, word string, idx int) string {
		upper := strings.ToUpper(word)
		if _, ok := initialisms[upper]; ok {
			return res + upper
		}
		word = strings.ToLower(word)
		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		return res + string(runes)
	})
}

type jointer func(result, word string, index int) string

func rewords(s string, fn jointer) string {
	words := SplitToWords(s)
	ret := ""

	for i, word := range words {
		ret = fn(ret, word, i)
	}
	return ret
}

var initialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
	"QOS":   true,
}
