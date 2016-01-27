package intseq

import (
	"bytes"
	"testing"
)

func Test(t *testing.T) {
	var vals = []uint64{1, 2, 3, 4, 5, 10, 2048, 2050, 2060}
	d := NewIntSeqFromSlice(vals)
	t.Log("d=", d)
	buf := bytes.NewBuffer([]byte{})
	d.WriteTo(buf)
	d2 := NewIntSeq()
	bufr := bytes.NewReader(buf.Bytes())
	d2.ReadFrom(bufr)
	if d.Count() != d2.Count() {
		t.Error("length mismatch", d, d2)
	}
	d.Each(func(v uint64) bool {
		if !d2.Contains(v) {
			t.Error(v, "is not in", d2)
		}
		return false
	})
}
