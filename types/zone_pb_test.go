package types

import (
	"testing"
)

func TestZonePb01(t *testing.T) {
	a := []*RecordPb{
		{
			Name: "test.com",
			Type: RecordPb_A,
			Ttl:  300,
			Value: []*RecordPb_Value{
				{
					Soa: &RecordPb_Value_SOA{
						Mname:   "test.com.",
						Nname:   "dns-admin.test.com.",
						Serial:  305419896,
						Refresh: 1193046,
						Retry:   624485,
						Expire:  4913,
						Minttl:  389333,
					},
				},
				{
					Record: []string{"1.1.1.1", "1.2.3.4"},
					View:   "dx",
				},
			},
		},
	}
	t.Logf("a : \n%s\n", a)
}
