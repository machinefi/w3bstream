package stringsx

import (
	"unicode"
	"unicode/utf8"
)

type LetterType int

const (
	LowerLetter LetterType = iota + 1
	UpperLetter
	Digit
	Other
)

func CheckLetterType(c rune) LetterType {
	switch {
	case unicode.IsLower(c):
		return LowerLetter
	case unicode.IsUpper(c):
		return UpperLetter
	case unicode.IsDigit(c):
		return Digit
	default:
		return Other
	}
}

func (t LetterType) IsLetter() bool { return t <= UpperLetter }
func (t LetterType) IsDigit() bool  { return t == Digit }
func (t LetterType) IsOther() bool  { return t == Other }

func SplitToWords(s string) []string {
	if !utf8.ValidString(s) {
		return []string{s}
	}

	words := make([]string, 0)
	runes := [][]rune{{rune(s[0])}}
	index := 1 // index for runes

	for i := 1; i < len(s); i++ {
		prev := CheckLetterType(runes[index-1][0])
		curr := CheckLetterType(rune(s[i]))
		if prev == curr {
			runes[index-1] = append(runes[index-1], rune(s[i]))
		} else {
			runes = append(runes, []rune{rune(s[i])})
			index++
		}
	}

	for i := 0; i < len(runes)-1; i++ {
		curr := CheckLetterType(runes[i][0])
		next := CheckLetterType(runes[i+1][0])
		if curr == UpperLetter && next == LowerLetter {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}

	for _, word := range runes {
		if len(word) == 0 {
			continue
		}
		kind := CheckLetterType(word[0])
		if kind == Other {
			continue
		}
		words = append(words, string(word))
	}

	return words
}
