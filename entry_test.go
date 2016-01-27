package intseq

import (
	"bytes"
	"testing"
)

func Test1(t *testing.T) {
	var vals = []uint64{1, 2, 3, 5, 6, 8, 10}
	d, _ := newIntSeqEntryFromSlice(vals)
	t.Log("d=", d)
	buf := bytes.NewBuffer([]byte{})
	d.writeTo(buf)
	d2 := newIntSeqEntry(0)
	bufr := bytes.NewReader(buf.Bytes())
	d2.readFrom(bufr)
	if d.initial != d2.initial {
		t.Error("initial not match")
	}
	if d.tail != d2.tail {
		t.Error("tail not match", d.tail, d2.tail)
	}
	d.each(func(v uint64) bool {
		if !d2.contains(v) {
			t.Error(v, "is not in", d2)
		}
		return false
	})
}

func TestNG(t *testing.T) {
	var vals = []uint64{1, 2, 3, 128, 1024}
	d, err := newIntSeqEntryFromSlice(vals)
	if err != nil {
		t.Log("mis", d, err)
		return
	}
	t.Fail()
}

func TestEach(t *testing.T) {
	var vals = []uint64{1, 2, 3, 4, 5}
	d, _ := newIntSeqEntryFromSlice(vals)
	var i = 0
	d.each(func(x uint64) bool {
		t.Log("x", x)
		if vals[i] != x {
			t.Error("mismatch", x, vals[i])
		}
		i++
		return false
	})
}
