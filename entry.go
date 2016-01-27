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

// IntSeqEntry is a entry
type IntSeqEntry struct {
	Initial uint64  // first value
	tail    uint64  // tail value
	Data    []uint8 // data
}

func NewIntSeqEntry(v uint64) (res *IntSeqEntry) {
	res = new(IntSeqEntry)
	res.Initial = v
	res.tail = v
	res.Data = make([]uint8, 0)
	return
}

func NewIntSeqEntryFromSlice(vals []uint64) (res *IntSeqEntry, err error) {
	res = new(IntSeqEntry)
	res.Initial = vals[0]
	res.tail = vals[0]
	res.Data = make([]uint8, 0)
	for _, v := range vals[1:] {
		if err = res.Add(v); err != nil {
			return nil, err
		}
	}
	return
}

func (self *IntSeqEntry) Add(v uint64) (e error) {
	if len(self.Data) > maxentry {
		return fmt.Errorf("too large")
	}
	diff := v - self.tail
	if 0 < diff && diff < 256 {
		self.Data = append(self.Data, uint8(diff))
		self.tail = v
	} else {
		return fmt.Errorf("range mismatch: %d", diff)
	}
	return
}

func (self *IntSeqEntry) Get(i uint) uint64 {
	cur := self.Initial
	for n := uint(0); n < i; {
		cur, n = self.Next(cur, n)
	}
	return cur
}

func (self *IntSeqEntry) Each(fn func(uint64) bool) bool {
	cur := self.Initial
	if fn(cur) {
		return true
	}
	for _, v := range self.Data {
		cur += uint64(v)
		if fn(cur) {
			return true
		}
	}
	return false
}

func (self *IntSeqEntry) Map(fn func(uint64) interface{}) []interface{} {
	var v = []interface{}{}
	self.Each(func(x uint64) bool {
		v = append(v, fn(x))
		return false
	})
	return v
}

func (self *IntSeqEntry) MapToString(fn func(uint64) string) []string {
	var v = []string{}
	self.Each(func(x uint64) bool {
		v = append(v, fn(x))
		return false
	})
	return v
}
func (self *IntSeqEntry) Next(v uint64, i uint) (uint64, uint) {
	return v + uint64(self.Data[i]), i + uint(1)
}

func (self *IntSeqEntry) Contains(v uint64) (ret bool) {
	ret = false
	self.Each(func(x uint64) bool {
		if x == v {
			ret = true
		}
		return ret
	})
	return
}

func (self *IntSeqEntry) Count() uint64 {
	return uint64(1) + uint64(len(self.Data))
}

func (self *IntSeqEntry) WriteTo(out io.Writer) (err error) {
	if err = binary.Write(out, binary.BigEndian, self.Initial); err != nil {
		log.Println("write initial", self.Initial, err)
		return
	}
	var l uint32 = uint32(len(self.Data))
	if err = binary.Write(out, binary.BigEndian, l); err != nil {
		log.Println("write length", l, err)
		return
	}
	if err = binary.Write(out, binary.BigEndian, self.Data); err != nil {
		log.Println("write data", len(self.Data), self.Data, err)
	}
	return
}

func (self *IntSeqEntry) ReadFrom(in io.Reader) (err error) {
	if err = binary.Read(in, binary.BigEndian, &self.Initial); err != nil {
		return
	}
	var l uint32
	if err = binary.Read(in, binary.BigEndian, &l); err != nil {
		log.Println("read length", l, err)
		return
	}
	self.Data = make([]byte, l)
	if err = binary.Read(in, binary.BigEndian, &self.Data); err != nil {
		log.Println("read data", len(self.Data), self.Data, err)
		return
	}
	self.tail = self.Initial
	for _, v := range self.Data {
		self.tail += uint64(v)
	}
	return
}

func (self *IntSeqEntry) String() string {
	return strings.Join(self.MapToString(func(v uint64) string {
		return fmt.Sprintf("%d", v)
	}), ",")
}
