package re

import (
	"regexp"
	"regexp/syntax"
	"unicode"
)

var (
	digits         = []rune("09")
	lowercaseAlpha = []rune("az")
	uppercaseAlpha = []rune("AZ")
	alpha          = append(uppercaseAlpha, lowercaseAlpha...)
	alphanum       = append(alpha, digits...)
	whitespaces    = []rune("  \t\t\n\n")
)

func toSyntax(args ...interface{}) []*syntax.Regexp {
	res := make([]*syntax.Regexp, len(args))
	for i, arg := range args {
		switch arg.(type) {
		case string:
			res[i] = Then(arg.(string))
		case *syntax.Regexp:
			res[i] = arg.(*syntax.Regexp)
		default:
			panic(arg)
		}
	}
	return res
}

// appendRange returns the result of appending the range lo-hi to the class r.
// source: https://golang.org/src/regexp/syntax/parse.go
func appendRange(r []rune, lo, hi rune) []rune {
	// Expand last range or next to last range if it overlaps or abuts.
	// Checking two ranges helps when appending case-folded
	// alphabets, so that one range can be expanding A-Z and the
	// other expanding a-z.
	n := len(r)
	for i := 2; i <= 4; i += 2 { // twice, using i=2, i=4
		if n >= i {
			rlo, rhi := r[n-i], r[n-i+1]
			if lo <= rhi+1 && rlo <= hi+1 {
				if lo < rlo {
					r[n-i] = lo
				}
				if hi > rhi {
					r[n-i+1] = hi
				}
				return r
			}
		}
	}

	return append(r, lo, hi)
}

// appendNegatedClass returns the result of appending the negation of the class x to the class r.
// source: https://golang.org/src/regexp/syntax/parse.go
func appendNegatedClass(r []rune, x []rune) []rune {
	nextLo := '\u0000'
	for i := 0; i < len(x); i += 2 {
		lo, hi := x[i], x[i+1]
		if nextLo <= lo-1 {
			r = appendRange(r, nextLo, lo-1)
		}
		nextLo = hi + 1
	}
	if nextLo <= unicode.MaxRune {
		r = appendRange(r, nextLo, unicode.MaxRune)
	}
	return r
}

// appendClass returns the result of appending the class x to the class r.
// It assume x is clean.
// source: https://golang.org/src/regexp/syntax/parse.go
func appendClass(r []rune, x []rune) []rune {
	for i := 0; i < len(x); i += 2 {
		r = appendRange(r, x[i], x[i+1])
	}
	return r
}

// negateClass overwrites r and returns r's negation.
// It assumes the class r is already clean.
// taken from https://golang.org/src/regexp/syntax/parse.go
func negateClass(r []rune) []rune {
	nextLo := '\u0000' // lo end of next class to add
	w := 0             // write index
	for i := 0; i < len(r); i += 2 {
		lo, hi := r[i], r[i+1]
		if nextLo <= lo-1 {
			r[w] = nextLo
			r[w+1] = lo - 1
			w += 2
		}
		nextLo = hi + 1
	}
	r = r[:w]
	if nextLo <= unicode.MaxRune {
		// It's possible for the negation to have one more
		// range - this one - than the original class, so use append.
		r = append(r, nextLo, unicode.MaxRune)
	}
	return r
}

var (
	StartOfLine = &syntax.Regexp{
		Op: syntax.OpBeginLine,
	}
	EndOfLine = &syntax.Regexp{
		Op: syntax.OpEndLine,
	}
	StartOfText = &syntax.Regexp{
		Op: syntax.OpBeginText,
	}
	EndOfText = &syntax.Regexp{
		Op: syntax.OpEndText,
	}
	Digit = &syntax.Regexp{
		Op:   syntax.OpCharClass,
		Rune: digits,
	}
	Period = &syntax.Regexp{
		Op:   syntax.OpLiteral,
		Rune: []rune{'.'},
	}
	Digits = &syntax.Regexp{
		Op:  syntax.OpPlus,
		Sub: []*syntax.Regexp{Digit},
	}
	Lowercase = &syntax.Regexp{
		Op:   syntax.OpCharClass,
		Rune: lowercaseAlpha,
	}
	Uppercase = &syntax.Regexp{
		Op:   syntax.OpCharClass,
		Rune: uppercaseAlpha,
	}
	Alpha = &syntax.Regexp{
		Op:   syntax.OpCharClass,
		Rune: alpha,
	}
	Alphanum = &syntax.Regexp{
		Op:   syntax.OpCharClass,
		Rune: alphanum,
	}
	Anything = &syntax.Regexp{
		Op: syntax.OpAnyChar,
	}
	Word = &syntax.Regexp{
		Op:  syntax.OpPlus,
		Sub: []*syntax.Regexp{Alphanum},
	}
	Newline = &syntax.Regexp{
		Op:   syntax.OpLiteral,
		Rune: []rune("\n"),
	}
	Whitespace = &syntax.Regexp{
		Op: syntax.OpPlus,
		Sub: []*syntax.Regexp{
			&syntax.Regexp{
				Op:   syntax.OpCharClass,
				Rune: whitespaces,
			},
		},
	}
)

func AnythingBut(args ...rune) *syntax.Regexp {
	if n := len(args); n%2 == 1 {
		args = append(args, args[n-1])
	}
	neg := negateClass(appendClass(nil, args))
	return &syntax.Regexp{
		Op:   syntax.OpCharClass,
		Rune: neg,
	}
}

func Range(rng ...rune) *syntax.Regexp {
	return &syntax.Regexp{
		Op:   syntax.OpCharClass,
		Rune: appendClass(nil, rng),
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

func Maybe(sub ...interface{}) *syntax.Regexp {
	return &syntax.Regexp{
		Op: syntax.OpQuest,
		Sub: []*syntax.Regexp{
			&syntax.Regexp{
				Op:  syntax.OpConcat,
				Sub: toSyntax(sub...),
			},
		},
	}
}

func AtLeastOne(sub ...*syntax.Regexp) *syntax.Regexp {
	return &syntax.Regexp{
		Op:  syntax.OpPlus,
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

func Sequence(subs ...*syntax.Regexp) *syntax.Regexp {
	return &syntax.Regexp{
		Op:  syntax.OpConcat,
		Sub: subs,
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
