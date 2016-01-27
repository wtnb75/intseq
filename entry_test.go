package intseq

import (
	"bytes"
	"testing"
)

func Test1(t *testing.T) {
	var vals = []uint64{1, 2, 3, 5, 6, 8, 10}
	d, _ := NewIntSeqEntryFromSlice(vals)
	t.Log("d=", d)
	buf := bytes.NewBuffer([]byte{})
	d.WriteTo(buf)
	d2 := NewIntSeqEntry(0)
	bufr := bytes.NewReader(buf.Bytes())
	d2.ReadFrom(bufr)
	if d.Initial != d2.Initial {
		t.Error("initial not match")
	}
	if d.tail != d2.tail {
		t.Error("tail not match", d.tail, d2.tail)
	}
}

func TestNG(t *testing.T) {
	var vals = []uint64{1, 2, 3, 128, 1024}
	d, err := NewIntSeqEntryFromSlice(vals)
	if err != nil {
		t.Log("mis", d, err)
		return
	}
	t.Fail()
}

func TestEach(t *testing.T) {
	var vals = []uint64{1, 2, 3, 4, 5}
	d, _ := NewIntSeqEntryFromSlice(vals)
	var i = 0
	d.Each(func(x uint64) bool {
		t.Log("x", x)
		if vals[i] != x {
			t.Error("mismatch", x, vals[i])
		}
		i++
		return false
	})
}
