package easy

import (
	"regexp"
	"regexp/syntax"
)

var (
	digits         = []rune("0123456789")
	lowercaseAlpha = []rune("abcdefghijklmnopqrstuvxywz")
	uppercaseAlpha = []rune("ABCDEFGHIJKLMNOPQRSTUVXYWZ")
	alpha          = append(uppercaseAlpha, lowercaseAlpha...)
	alphanum       = append(alpha, digits...)
)

func Digit() *syntax.Regexp {
	return &syntax.Regexp{
		Op:   syntax.OpCharClass,
		Rune: digits,
	}
}

func Period() *syntax.Regexp {
	return &syntax.Regexp{
		Op:   syntax.OpLiteral,
		Rune: []rune{'.'},
	}
}

func Digits() *syntax.Regexp {
	return &syntax.Regexp{
		Op:  syntax.OpPlus,
		Sub: []*syntax.Regexp{Digit()},
	}
}

func Alpha() *syntax.Regexp {
	return &syntax.Regexp{
		Op:   syntax.OpCharClass,
		Rune: alpha,
	}
}

func Alphanum() *syntax.Regexp {
	return &syntax.Regexp{
		Op:   syntax.OpCharClass,
		Rune: alphanum,
	}
}

func Word() *syntax.Regexp {
	return &syntax.Regexp{
		Op:  syntax.OpPlus,
		Sub: []*syntax.Regexp{Alphanum()},
	}
}

func Then(match string) *syntax.Regexp {
	return &syntax.Regexp{
		Op:   syntax.OpLiteral,
		Rune: []rune(match),
	}
}

func Repeat(times int, sub ...*syntax.Regexp) *syntax.Regexp {
	return &syntax.Regexp{
		Op:  syntax.OpRepeat,
		Min: times,
		Max: times,
		Sub: sub,
	}
}

func Max(times int, sub ...*syntax.Regexp) *syntax.Regexp {
	return &syntax.Regexp{
		Op:  syntax.OpRepeat,
		Max: times,
		Sub: sub,
	}
}

func Min(times int, sub ...*syntax.Regexp) *syntax.Regexp {
	return &syntax.Regexp{
		Op:  syntax.OpRepeat,
		Min: times,
		Sub: sub,
	}
}

func Or(sub ...*syntax.Regexp) *syntax.Regexp {
	return &syntax.Regexp{
		Op:  syntax.OpAlternate,
		Sub: sub,
	}
}

/*
Groups are currently only supported at the topmost level of Regex
*/
func Group(name string, sub ...*syntax.Regexp) *syntax.Regexp {
	return &syntax.Regexp{
		Op:   syntax.OpCapture,
		Sub:  sub,
		Name: name,
	}
}

func CompileRegex(subs ...*syntax.Regexp) (*regexp.Regexp, error) {
	capIdx := 0
	for _, sub := range subs {
		if sub.Op == syntax.OpCapture {
			sub.Cap = capIdx
			capIdx++
		}
	}
	re := &syntax.Regexp{
		Op:  syntax.OpConcat,
		Sub: subs,
	}
	re = re.Simplify()
	return regexp.Compile(re.String()) // see https://github.com/golang/go/issues/18888]
}

func Regex(subs ...*syntax.Regexp) *regexp.Regexp {

	re, err := CompileRegex(subs...)
	if err != nil {
		panic(err)
	}
	return re
}
