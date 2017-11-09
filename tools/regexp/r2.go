package main

import (
	//"github.com/umisama/go-regexpcache"
	"fmt"
	"strings"
	"sort"
)

type SlagSlice []string

func (s SlagSlice) Len() int {
    return len(s)
}
func (s SlagSlice) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s SlagSlice) Less(i, j int) bool {
	iwL := strings.Count(s[i], `\w+`)
	idL := strings.Count(s[i], `\.`)
	jwL := strings.Count(s[j], `\w+`)
	jdL := strings.Count(s[j], `\.`)
	iL :=  iwL +  idL * 2 + (idL + 1 - iwL) * 3
	jL := jwL +  jdL * 2 + (jdL + 1 - jwL) * 3
    return iL > jL
}

func main() {

	//str := `asdf.a`
	fruits := []string{
			`\w+\.a`,
			`\w+\.blaaklfjsdsf`,			
			`\w+\.a\.\w+\.a`,
			`\w+\.a\.asdfsdf`,
			`\w+\.a\.\w+`,
			`\w+\.a\.g`,
			}
	sort.Sort(SlagSlice(fruits))
    fmt.Println(fruits)

}
