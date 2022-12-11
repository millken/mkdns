package zone

import (
	"sort"
	"testing"

	"github.com/millken/mkdns/dns"
	"github.com/stretchr/testify/require"
)

func TestSortKeyList(t *testing.T) {
	a := KeyList{
		Key{"*.123.a", dns.TypeA},
		Key{"*.a", dns.TypeA},
		Key{"*.a3", dns.TypeA},
		Key{"*.b", dns.TypeA},
		Key{"*.a.a.b", dns.TypeA},
	}
	sort.Sort(a)
	t.Log(a)
}

func TestDebugRecord(t *testing.T) {
	require := require.New(t)
	z := New("c.com")
	z.Add("*", &dns.Record{Name: "*", Type: dns.TypeA, TTL: 0, Value: []dns.RecordValue{{Data: []string{"aa"}}}})
	z.Add("b", &dns.Record{Name: "b", Type: dns.TypeA, TTL: 0, Value: []dns.RecordValue{{Data: []string{"a0"}}}})
	z.Add("a.b", &dns.Record{Name: "a.b", Type: dns.TypeA, TTL: 0, Value: []dns.RecordValue{{Data: []string{"a1"}}}})
	z.Add("*.b", &dns.Record{Name: "*.b", Type: dns.TypeA, TTL: 0, Value: []dns.RecordValue{{Data: []string{"a2"}}}})
	r, ok := z.Lookup(Key{"aa", dns.TypeA})
	require.True(ok)
	require.Equal(r.Value[0].Data[0], "aa")
	r, ok = z.Lookup(Key{"b", dns.TypeA})
	require.True(ok)
	require.Equal(r.Value[0].Data[0], "a0")
	r, ok = z.Lookup(Key{"a.b", dns.TypeA})
	require.True(ok)
	require.Equal(r.Value[0].Data[0], "a1")
	r, ok = z.Lookup(Key{"a3.b", dns.TypeA})
	require.True(ok)
	require.Equal(r.Value[0].Data[0], "a2")
	r, ok = z.Lookup(Key{"a3.b3", dns.TypeA})
	require.True(ok)
	require.Equal(r.Value[0].Data[0], "aa")
}

func BenchmarkZoneAdd(b *testing.B) {
	z := New("test.com")
	for i := 0; i < b.N; i++ {
		z.Add("a.b.c.d", &dns.Record{Name: "a.b.c.d", Type: dns.TypeA, TTL: 0, Value: []dns.RecordValue{{Data: []string{"aa"}}}})
	}
}

func BenchmarkZoneLookup(b *testing.B) {
	z := New("test.com")
	z.Add("a.b.c.d", &dns.Record{Name: "a.b.c.d", Type: dns.TypeA, TTL: 0, Value: []dns.RecordValue{{Data: []string{"aa"}}}})
	z.Add("*.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z", &dns.Record{Name: "a.b.c.d", Type: dns.TypeA, TTL: 0, Value: []dns.RecordValue{{Data: []string{"aa"}}}})

	for i := 0; i < b.N; i++ {
		z.Lookup(Key{"a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z", dns.TypeA})
	}
}
