package re

import (
	"regexp"
	"testing"
)

func TestRegex(t *testing.T) {
	type testcase struct {
		str       string
		doesMatch bool
	}

	testdata := []struct {
		re        *regexp.Regexp
		ref       *regexp.Regexp
		testcases []testcase
	}{
		{
			Regex(Digit()),
			regexp.MustCompile("[0-9]"),
			[]testcase{
				{"0", true},
				{"1", true},
				{"9", true},
				{"9.", true},
				{"-9", true},
				{"-999.2", true},
				{" asdasda s -999.2\n asdasdsad", true},
				{"", false},
				{"asdasd", false},
				{"digit", false},
			},
		},
		{
			Regex(Digits(), Then("."), Digits()),
			regexp.MustCompile("[0-9]+\\.[0-9]+"),
			[]testcase{
				{"1320.0", true},
				{"2.1", true},
				{"0.", false},
				{"0", false},
				{"   0.0", true},
				{"here is > 0.0", true},
				{"here is not", false},
			},
		},
		{
			Regex(Digits(), Period(), Digits()),
			regexp.MustCompile("[0-9]+\\.[0-9]+"),
			[]testcase{
				{"1320.0", true},
				{"2.1", true},
				{"0.", false},
				{"0", false},
				{"   0.0", true},
				{"here is > 0.0", true},
				{"here is not", false},
			},
		},
		{
			Regex(Word()),
			regexp.MustCompile("\\w+"),
			[]testcase{
				{"the quick brown fox jumps over the lazy dog", true},
				{"the", true},
				{"quick", true},
				{"brown", true},
				{"fox", true},
				{"jumps", true},
				{"over", true},
				{"lazy", true},
				{"dog", true},
				{"a", true},
				{"zzzz", true},
				{"", false},
				{"abc999", true},
				{"999", true},
				{"----", false},
				{"", false},
			},
		},
		{
			Regex(Then("$"), Digits(), Then("USD")),
			regexp.MustCompile("\\$\\d+USD"),
			[]testcase{
				{"$123USD", true},
				{"$0USD", true},
				{"$0.", false},
				{"$USD", false},
				{"abc", false},
				{"  $99999999999999USD    ", true},
				{"  $999 USD ", false},
			},
		},
		{
			Regex(Repeat(3, Then("a")), Then("b")),
			regexp.MustCompile("aaab"),
			[]testcase{
				{"aaab", true},
				{"aab", false},
				{" aaaab", true},
				{" aaaaabbbbb ", true},
				{"aabaab", false},
				{"abaabaaabbb", true},
			},
		},
		{
			Regex(Max(3, Then("a")), Then("b")),
			regexp.MustCompile("a{0,3}b"),
			[]testcase{
				{"aaab", true},
				{"aab", true},
				{" aaaab", true},
				{" aaaaabbbbb ", true},
				{"aabaab", true},
				{"abaabaaabbb", true},
				{"bbb", true},
			},
		},
		{
			Regex(Min(3, Then("a")), Then("b")),
			regexp.MustCompile("a{3,}b"),
			[]testcase{
				{"", false},
				{"b", false},
				{"aab", false},
				{"aaab", true},
				{"aaaab", true},
				{"aaaaab", true},
				{"aaaaabb", true},
			},
		},
		{
			Regex(Or(Digit(), Alpha())),
			regexp.MustCompile("[A-Za-z]|[0-9]"),
			[]testcase{
				{"", false},
				{"b", true},
				{"1", true},
				{"z", true},
				{"-", false},
				{" ", false},
				{"X", true},
				{"A", true},
			},
		},
		{
			Regex(Or(Then("aaa"), Then("bbb"))),
			regexp.MustCompile("(aaa)|(bbb)"),
			[]testcase{
				{"", false},
				{"aaa", true},
				{"bbb", true},
				{"z", false},
				{"AAA", false},
			},
		},
		{
			Regex(StartOfLine(), Digits(), EndOfLine()),
			regexp.MustCompilePOSIX("^[0-9]+$"),
			[]testcase{
				{"", false},
				{"1234567890", true},
				{"1", true},
				{" 123", false},
				{"123 ", false},
				{"abc\n123\nxyz", true},
				{"1x1", false},
			},
		},
		{
			Regex(StartOfText(), Digits(), EndOfText()),
			regexp.MustCompile("\\A[0-9]+\\z"),
			[]testcase{
				{"", false},
				{"1234567890", true},
				{"1", true},
				{" 123", false},
				{"123 ", false},
				{"abc\n123\nxyz", false},
				{"123\n", false},
				{"\n123", false},
				{"1x1", false},
			},
		},
		{
			Regex(
				Group("dividend", Digits()),
				Then("/"),
				Group("divisor", Digits()),
			),
			regexp.MustCompile("(\\d+)/(\\d+)"),
			[]testcase{
				{"", false},
				{"1", false},
				{"/", false},
				{"1/", false},
				{"1/2", true},
				{"99/9", true},
				{"91231239/1231231239", true},
				{"-91231239/1231231239", true},
				{"+91231239/1231231239", true},
				{" +91231239/1231231239", true},
				{" +/2", false},
				{"-/9", false},
			},
		},
		{
			Regex(
				Then("a"),
				Maybe(Digit()),
			),
			regexp.MustCompile("a[0-9]?"),
			[]testcase{
				{"", false},
				{"a", true},
				{"a1", true},
				{"a2", true},
				{"a0", true},
				{"11", false},
				{"1", false},
			},
		},
	}

	for _, td := range testdata {
		t.Run(td.re.String(), func(t *testing.T) {
			for _, tc := range td.testcases {
				t.Run(tc.str, func(t *testing.T) {
					if td.ref.MatchString(tc.str) != tc.doesMatch {
						t.Fatal(tc.str, tc.doesMatch)
					}

					if td.re.MatchString(tc.str) != tc.doesMatch {
						t.Fatal(tc.str, tc.doesMatch)
					}
				})
			}
		})
	}
}

func TestRegexGroups(t *testing.T) {
	type testcase struct {
		str     string
		matches []string
	}

	testdata := []struct {
		re        *regexp.Regexp
		ref       *regexp.Regexp
		testcases []testcase
	}{
		{
			Regex(
				Group("dividend", Digits()),
				Then("/"),
				Group("divisor", Digits()),
			),
			regexp.MustCompile("(\\d+)\\.(\\d+)"),
			[]testcase{
				{"", nil},
				{"1", nil},
				{"/", nil},
				{"1/", nil},
				{"1/2", []string{"1/2", "1", "2"}},
				{"99/9", []string{"99/9", "99", "9"}},
				{"91231239/1231231239", []string{"91231239/1231231239", "91231239", "1231231239"}},
				{"-91231239/1231231239", []string{"91231239/1231231239", "91231239", "1231231239"}},
				{"+91231239/1231231239", []string{"91231239/1231231239", "91231239", "1231231239"}},
				{" +321/765", []string{"321/765", "321", "765"}},
				{" +/2", nil},
				{"-/9", nil},
			},
		},
		{
			Regex(
				Group("user", Word()),
				Then("@"),
				Group("domain", Word()),
			),
			regexp.MustCompile("(\\w+)@(\\w+)"),
			[]testcase{
				{"foo-bar.com", nil},
				{"foo@bar", []string{"foo@bar", "foo", "bar"}},
				{"foo@bar.com", []string{"foo@bar", "foo", "bar"}},
			},
		},
	}

	for _, td := range testdata {
		t.Run(td.re.String(), func(t *testing.T) {
			for _, tc := range td.testcases {
				t.Run(tc.str, func(t *testing.T) {
					submatches := td.re.FindSubmatch([]byte(tc.str))
					if len(submatches) != len(tc.matches) {
						t.Fatal(submatches, len(submatches))
					}

					for i := range submatches {
						if string(submatches[i]) != tc.matches[i] {
							t.Fatal(string(submatches[i]), tc.matches[i])
						}
					}
				})
			}
		})
	}
}
