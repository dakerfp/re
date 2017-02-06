Golang Regular Expressions for Humans
=====================================

This module is inspired in a similar project for [javascript](https://github.com/VerbalExpressions/JSVerbalExpressions).


Example:

```golang
import (
	. "github.com/dakerfp/re"
)

re := Regex(
	Group("dividend", Digits),
	Then("/"),
	Group("divisor", Digits),
)

m = re.FindSubmatch("4/3")
fmt.Println(m[1]) // > 4
fmt.Println(m[2]) // > 3
```

the equivalent regexp would be:

```golang
regexp.MustCompile("(\\d+)/(\\d+)")
```

which is far more cryptic.

Another good example is the following regex (limited) to parse URLs:

```golang
import . "github.com/dakerfp/re"

re := re.Regex(
	StartOfLine,
	Then("http"),
	Maybe("s"),
	Then("://"),
	Maybe("www"),
	AtLeastOne(AnythingBut(' ')),
	EndOfLine,
)
```

which corresponds to

```golang
regexp.MustCompile("^http[s]?://(www)?[^\\ ]+$"),
```

Contributing
------------

Fill issues for discussions and request, pull request for
fixes and new features.
License is Apache 2.0.


TODO List
---------

* Create Docs
* Validate API
* Add more tests