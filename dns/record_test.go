package dns

import (
	"testing"

	"golang.org/x/net/publicsuffix"
)

func BenchmarkPublicSuffix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		publicsuffix.PublicSuffix("a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z")
	}
}
