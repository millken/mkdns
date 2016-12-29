package main

import (
	"github.com/umisama/go-regexpcache"
	"fmt"
)

func main() {
	cases := []struct {
		exp          string
		str          string
	}{
		{`\w+.a.\w+`, "aaaaaa.f.a"},
		{"\\w+.a", "aaaa.a"},
		{"\\w+.a", "aaaa.f"},
		{`\w+\.a`, "aaaa.a"},
	}
	for _, c := range cases {
		match := regexpcache.MustCompile(c.exp).MatchString(c.str)
		fmt.Printf("str: %s, exp:%s, match: %+v\n", c.str, c.exp, match)
	}

}
