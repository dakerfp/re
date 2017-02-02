package easy

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
			regexp.MustCompile("\\d"),
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
			regexp.MustCompile("\\d+\\.\\d+"),
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
			regexp.MustCompile("\\d+\\.\\d+"),
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
				/// {"zzzz", true},
				{"", false},
				{"999", false},
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
		// {
		// 	Regex(BeginText(), Then("abc")), // XXX
		// 	regexp.MustCompile("$abc"),
		// 	[]testcase{
		// 		{"abc", true},
		// 		{"abcd", true},
		// 		{"abcd\nabcd", false},
		// 		{"\nabcd", false},
		// 		{"aabcd", false},
		// 		{"ab", false},
		// 		{" abcd", false},
		// 		{"", false},
		// 		{"", true}, // XXX
		// 	},
		// },
		// {
		// 	Regex(BeginLine(), Then("abc")), // XXX
		// 	regexp.MustCompile("$abc"),
		// 	[]testcase{
		// 		{"abc", true},
		// 		{"abcd", true},
		// 		{"acv\nabcd", true},
		// 		{"\nabcd", true},
		// 		{"aabcd", false},
		// 		{"ab", false},
		// 		{" abcd", false},
		// 		{"  \nabcd", false},
		// 		{"", false},
		// 	},
		// },
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
	return //

	type testcase struct {
		str       string
		doesMatch bool
		matches   []string
	}

	testdata := []struct {
		re        *regexp.Regexp
		testcases []testcase
	}{
		{
			Regex(
				Group("dividend", Digits()),
				Then("/"),
				Group("divisor", Digits()),
			),
			[]testcase{
				{"", false, nil},
				{"1", false, nil},
				{"/", false, nil},
				{"1/", false, nil},
				{"1/2", true, []string{"1", "2"}},
				{"99/9", true, []string{"99", "9"}},
				{"91231239/1231231239", true, []string{"91231239", "1231231239"}},
				{"-91231239/1231231239", true, []string{"91231239", "1231231239"}},
				{"+91231239/1231231239", true, []string{"91231239", "1231231239"}},
				{" +91231239/1231231239", true, []string{"91231239", "1231231239"}},
				{" +/2", false, nil},
				{"-/9", false, nil},
			},
		},
	}

	for _, td := range testdata {
		t.Run(td.re.String(), func(t *testing.T) {
			t.Log(td.re.String)
			for _, tc := range td.testcases {
				t.Run(tc.str, func(t *testing.T) {
					if td.re.MatchString(tc.str) != tc.doesMatch {
						t.Fail()
					}

					if !tc.doesMatch {
						return
					}

					submatches := td.re.FindSubmatch([]byte(tc.str))
					if len(submatches) == 0 || len(submatches) != len(tc.matches) {
						t.Fatal(submatches)
					}

					for i := range submatches {
						if string(submatches[i]) != tc.matches[i] {
							t.Fatal(submatches[i], tc.matches[i])
						}
					}
				})
			}
		})
	}
}
