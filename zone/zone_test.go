package zone

import (
	"sort"
	"testing"

	"github.com/millken/mkdns/internal/wire"
	"github.com/stretchr/testify/require"
)

func TestSortKeyList(t *testing.T) {
	a := KeyList{
		Key{"*.123.a", wire.TypeA},
		Key{"*.a", wire.TypeA},
		Key{"*.a3", wire.TypeA},
		Key{"*.b", wire.TypeA},
		Key{"*.a.a.b", wire.TypeA},
	}
	sort.Sort(a)
	t.Log(a)
}

func TestDebugRecord(t *testing.T) {
	require := require.New(t)
	z := New("c.com")
	z.Add("*", &wire.Record{Name: "*", Type: wire.TypeA, TTL: 0, Value: []wire.RecordValue{{Data: []string{"aa"}}}})
	z.Add("b", &wire.Record{Name: "b", Type: wire.TypeA, TTL: 0, Value: []wire.RecordValue{{Data: []string{"a0"}}}})
	z.Add("a.b", &wire.Record{Name: "a.b", Type: wire.TypeA, TTL: 0, Value: []wire.RecordValue{{Data: []string{"a1"}}}})
	z.Add("*.b", &wire.Record{Name: "*.b", Type: wire.TypeA, TTL: 0, Value: []wire.RecordValue{{Data: []string{"a2"}}}})
	r, ok := z.Lookup(Key{"aa", wire.TypeA})
	require.True(ok)
	require.Equal(r.Value[0].Data[0], "aa")
	r, ok = z.Lookup(Key{"b", wire.TypeA})
	require.True(ok)
	require.Equal(r.Value[0].Data[0], "a0")
	r, ok = z.Lookup(Key{"a.b", wire.TypeA})
	require.True(ok)
	require.Equal(r.Value[0].Data[0], "a1")
	r, ok = z.Lookup(Key{"a3.b", wire.TypeA})
	require.True(ok)
	require.Equal(r.Value[0].Data[0], "a2")
	r, ok = z.Lookup(Key{"a3.b3", wire.TypeA})
	require.True(ok)
	require.Equal(r.Value[0].Data[0], "aa")
}

func BenchmarkZoneAdd(b *testing.B) {
	z := New("test.com")
	for i := 0; i < b.N; i++ {
		z.Add("a.b.c.d", &wire.Record{Name: "a.b.c.d", Type: wire.TypeA, TTL: 0, Value: []wire.RecordValue{{Data: []string{"aa"}}}})
	}
}

func BenchmarkZoneLookup(b *testing.B) {
	z := New("test.com")
	z.Add("a.b.c.d", &wire.Record{Name: "a.b.c.d", Type: wire.TypeA, TTL: 0, Value: []wire.RecordValue{{Data: []string{"aa"}}}})
	z.Add("*.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z", &wire.Record{Name: "a.b.c.d", Type: wire.TypeA, TTL: 0, Value: []wire.RecordValue{{Data: []string{"aa"}}}})

	for i := 0; i < b.N; i++ {
		z.Lookup(Key{"a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z", wire.TypeA})
	}
}
