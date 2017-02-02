Easy Regular Expressions for Golang
===================================

This module is inspired on a similar project for [[javascript|https://github.com/VerbalExpressions/JSVerbalExpressions]].


Example:

```golang
import (
	. "github.com/dakerfp/re"
)
re := Regex(
	Group("dividend", Digits()),
	Then("/"),
	Group("divisor", Digits()),
)

m = re.FindSubmatch("4/3")
fmt.Println(m[1]) // > 4
fmt.Println(m[2]) // > 4
```
