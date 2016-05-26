package types

import (
	"testing"
)

var (
	err error
	content = `
@	3600	SOA  {"mname": "ns1.test1.com.", "nname": "dns-admin.test1.com.", "serial": 305419896, "refresh": 1193046,"retry": 624485, "expire": 4913, "minttl": 38933}

@ 600 NS {"record":["ns1.test1.com.", "ns2.test2.com"]}
@	600   A	{"records" : [{"record": ["1.47.46.2", "1.2.3.3"]}]}
@	600   AAAA	{"records" : [{"record": ["2001:0db8:85a3:08d3:1319:8a2e:0370:7344"]}]}
@	600   MX	{"records" : [{"record": [{"value": "mxbiz2.qq.com", "mx": 5},{"value": "mxbiz1.qq.com.", "mx": 50}]}]}
view	600   A	{"type": 1, "records" : [{"record": ["1.47.46.2", "1.2.3.3"], "view": "any"}, {"record": ["1.2.3.4", "1.2.4.5"], "view": "dx"}]}
weight	600   A	{"type": 2, "records" : [{"record": ["1.47.46.2", "1.2.3.3"], "weight": 3}, {"record": ["1.2.3.4", "1.2.4.5"], "weight": 7}, {"record": ["7.7.7.7"], "weight": 10}]}
b   312  A  {"type": 3,  "records" : [{"record": ["1.47.46.2", "1.2.3.3"], "weight":3, "view":"any"},{"record": ["1.47.46.2", "1.2.3.3"], "weight":9}]}
geo-1 600   A  {"type": 4, "records" : [{"record": ["7.7.7.7"]}, {"record": ["1.7.6.2"], "continent": "asia"}, {"record": ["1.7.6.5"], "continent": "asia", "country": "cn"}, {"record": ["1.7.6.6"], "country": "cn"}, {"record": ["1.2.3.4", "1.2.4.5"], "country": "kr"}, {"record": ["1.1.1.1", "1.2.2.3", "1.1.1.2"], "continent": "north-america"}, {"record": ["1.1.1.3"], "country": "us"}]}

*     353   A  {"records" : [{"record": ["1.2.3.4"]}]}
*.bb     353   A  {"records" : [{"record": ["1.2.3.5"]}]}
aa.*     353   A  {"records" : [{"record": ["1.2.3.6"]}]}
*.*     353   A  {"records" : [{"record": ["1.2.3.7"]}]}
a*b*.c 353   A  {"records" : [{"record": ["1.2.3.8"]}]}
^aa.*$ 353   A  {"records" : [{"record": ["1.2.3.9"]}]}
bücher 22  A {"records" : [{"record": ["2.2.3.9"]}]}
中文 22  A {"records" : [{"record": ["2.2.3.0"]}]}
@ 600 TXT {"records": [{"record": "AaBbCcDdEeFf"}]}
c1 600 CNAME {"records": [{"record": ["aaa.aaaa.aaaa.aaaa.com"]}]}
`;
)

func TestParse01(t *testing.T) {
	z := NewZone()
	if err = z.ParseBody([]byte(content)); err != nil {
		t.Fatal("parse failed:", err)
	}
}

// Benchmark command
//	go test -bench=Find
//	BenchmarkFind 1000000       1440 ns/op
func BenchmarkParse(b *testing.B) {
	b.StopTimer()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		z := NewZone()
		if err = z.ParseBody([]byte(content)); err != nil {
			b.Fatal("parse failed:", err)
		}
	}
}

