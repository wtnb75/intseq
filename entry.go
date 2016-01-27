package intseq

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

const (
	maxentry int = 1024 * 1024
)

// intSeqEntry is a entry
type intSeqEntry struct {
	initial uint64  // first value
	tail    uint64  // tail value
	data    []uint8 // data
}

func newIntSeqEntry(v uint64) (res *intSeqEntry) {
	res = new(intSeqEntry)
	res.initial = v
	res.tail = v
	res.data = make([]uint8, 0)
	return
}

func newIntSeqEntryFromSlice(vals []uint64) (res *intSeqEntry, err error) {
	res = new(intSeqEntry)
	res.initial = vals[0]
	res.tail = vals[0]
	res.data = make([]uint8, 0)
	for _, v := range vals[1:] {
		if err = res.add(v); err != nil {
			return nil, err
		}
	}
	return
}

func (iseq *intSeqEntry) add(v uint64) (e error) {
	if len(iseq.data) > maxentry {
		return fmt.Errorf("too large")
	}
	diff := v - iseq.tail
	if 0 < diff && diff < 256 {
		iseq.data = append(iseq.data, uint8(diff))
		iseq.tail = v
	} else {
		return fmt.Errorf("range mismatch: %d", diff)
	}
	return
}

func (iseq *intSeqEntry) get(i uint) uint64 {
	cur := iseq.initial
	for n := uint(0); n < i; {
		cur, n = iseq.next(cur, n)
	}
	return cur
}

func (iseq *intSeqEntry) each(fn func(uint64) bool) bool {
	cur := iseq.initial
	if fn(cur) {
		return true
	}
	for _, v := range iseq.data {
		cur += uint64(v)
		if fn(cur) {
			return true
		}
	}
	return false
}

func (iseq *intSeqEntry) mapfn(fn func(uint64) interface{}) []interface{} {
	var v = []interface{}{}
	iseq.each(func(x uint64) bool {
		v = append(v, fn(x))
		return false
	})
	return v
}

func (iseq *intSeqEntry) mapToString(fn func(uint64) string) []string {
	var v = []string{}
	iseq.each(func(x uint64) bool {
		v = append(v, fn(x))
		return false
	})
	return v
}
func (iseq *intSeqEntry) next(v uint64, i uint) (uint64, uint) {
	return v + uint64(iseq.data[i]), i + uint(1)
}

func (iseq *intSeqEntry) contains(v uint64) (ret bool) {
	ret = false
	if !(iseq.initial <= v && v <= iseq.tail) {
		return
	}
	iseq.each(func(x uint64) bool {
		if x == v {
			ret = true
		}
		return ret
	})
	return
}

func (iseq *intSeqEntry) count() uint64 {
	return uint64(1) + uint64(len(iseq.data))
}

func (iseq *intSeqEntry) writeTo(out io.Writer) (err error) {
	if err = binary.Write(out, binary.BigEndian, iseq.initial); err != nil {
		log.Println("write initial", iseq.initial, err)
		return
	}
	var l = uint32(len(iseq.data))
	if err = binary.Write(out, binary.BigEndian, l); err != nil {
		log.Println("write length", l, err)
		return
	}
	if err = binary.Write(out, binary.BigEndian, iseq.data); err != nil {
		log.Println("write data", len(iseq.data), iseq.data, err)
	}
	return
}

func (iseq *intSeqEntry) readFrom(in io.Reader) (err error) {
	if err = binary.Read(in, binary.BigEndian, &iseq.initial); err != nil {
		if err != io.EOF {
			log.Println("read initial", err)
		}
		return
	}
	var l uint32
	if err = binary.Read(in, binary.BigEndian, &l); err != nil {
		log.Println("read length", l, err)
		return
	}
	iseq.data = make([]byte, l)
	if err = binary.Read(in, binary.BigEndian, &iseq.data); err != nil {
		log.Println("read data", len(iseq.data), iseq.data, err)
		return
	}
	iseq.tail = iseq.initial
	for _, v := range iseq.data {
		iseq.tail += uint64(v)
	}
	return
}

func (iseq *intSeqEntry) String() string {
	return strings.Join(iseq.mapToString(func(v uint64) string {
		return fmt.Sprintf("%d", v)
	}), ",")
}
