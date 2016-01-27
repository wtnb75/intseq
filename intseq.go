package intseq

import (
	"io"
	"strings"
)

type IntSeq struct {
	Entries []IntSeqEntry
}

func NewIntSeq() *IntSeq {
	res := new(IntSeq)
	res.Entries = make([]IntSeqEntry, 0)
	return res
}

func NewIntSeqFromSlice(v []uint64) *IntSeq {
	res := NewIntSeq()
	for _, e := range v {
		res.Add(e)
	}
	return res
}

func (self *IntSeq) Clear() {
	self.Entries = make([]IntSeqEntry, 0)
}

func (self *IntSeq) Add(v uint64) error {
	if len(self.Entries) == 0 {
		self.Entries = append(self.Entries, *NewIntSeqEntry(v))
		return nil
	}
	target := &self.Entries[len(self.Entries)-1] // last entry
	err := target.Add(v)
	if err != nil {
		self.Entries = append(self.Entries, *NewIntSeqEntry(v))
	}
	return nil
}

func (self *IntSeq) Count() uint64 {
	var r uint64 = 0
	for _, e := range self.Entries {
		r += e.Count()
	}
	return r
}

func (self *IntSeq) Contains(v uint64) bool {
	for _, e := range self.Entries {
		if e.Contains(v) {
			return true
		}
	}
	return false
}

func (self *IntSeq) Each(fn func(uint64) bool) {
	for _, e := range self.Entries {
		if e.Each(fn) {
			break
		}
	}
}

func (self *IntSeq) WriteTo(wr io.Writer) error {
	for _, v := range self.Entries {
		if err := v.WriteTo(wr); err != nil {
			return err
		}
	}
	return nil
}

func (self *IntSeq) ReadFrom(rd io.Reader) error {
	for {
		r := NewIntSeqEntry(0)
		if err := r.ReadFrom(rd); err != nil {
			return nil
		}
		self.Entries = append(self.Entries, *r)
	}
	return nil
}

func (self *IntSeq) String() string {
	var s = []string{}
	for _, v := range self.Entries {
		s = append(s, v.String())
	}
	return "[" + strings.Join(s, ",") + "]"
}
