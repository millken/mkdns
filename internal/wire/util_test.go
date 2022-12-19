package wire

import (
	"testing"
)

func TestEncodeDomain(t *testing.T) {
	var cases = []struct {
		Domain string
		QName  string
	}{
		{"phus.lu", "\x04phus\x02lu\x00"},
		{"splunk.phus.lu", "\x06splunk\x04phus\x02lu\x00"},
	}

	for _, c := range cases {
		if got, want := string(EncodeDomain(nil, c.Domain)), c.QName; got != want {
			t.Errorf("EncodeDomain(%v) error got=%#v want=%#v", c.Domain, got, want)
		}
	}
}

func BenchmarkEncodeDomain(b *testing.B) {
	dst := make([]byte, 0, 256)
	for i := 0; i < b.N; i++ {
		dst = EncodeDomain(dst[:0], "hk.phus.lu")
	}
}
