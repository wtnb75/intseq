package intseq

import (
	"io"
	"strings"
)

// IntSeq is Integer Sequence Array
type IntSeq struct {
	Entries []intSeqEntry
}

// NewIntSeq creates new IntSeq structure
func NewIntSeq() *IntSeq {
	res := new(IntSeq)
	res.Entries = make([]intSeqEntry, 0)
	return res
}

// NewIntSeqFromSlice creates new IntSeq structure from integer slice
func NewIntSeqFromSlice(v []uint64) *IntSeq {
	res := NewIntSeq()
	for _, e := range v {
		res.Add(e)
	}
	return res
}

// Clear cleans up structure
func (iseq *IntSeq) Clear() {
	iseq.Entries = make([]intSeqEntry, 0)
}

// Add push back value
func (iseq *IntSeq) Add(v uint64) error {
	if len(iseq.Entries) == 0 {
		iseq.Entries = append(iseq.Entries, *newIntSeqEntry(v))
		return nil
	}
	target := &iseq.Entries[len(iseq.Entries)-1] // last entry
	err := target.add(v)
	if err != nil {
		iseq.Entries = append(iseq.Entries, *newIntSeqEntry(v))
	}
	return nil
}

// Count returns length
func (iseq *IntSeq) Count() uint64 {
	var r = uint64(0)
	for _, e := range iseq.Entries {
		r += e.count()
	}
	return r
}

// Length is alias for Count
func (iseq *IntSeq) Length() uint64 {
	return iseq.Count()
}

// Contains check if value is in array
func (iseq *IntSeq) Contains(v uint64) bool {
	for _, e := range iseq.Entries {
		if e.contains(v) {
			return true
		}
	}
	return false
}

// Each traverse values
func (iseq *IntSeq) Each(fn func(uint64) bool) {
	for _, e := range iseq.Entries {
		if e.each(fn) {
			break
		}
	}
}

// WriteTo writes data into io.Writer
func (iseq *IntSeq) Write(wr io.Writer) error {
	for _, v := range iseq.Entries {
		if err := v.writeTo(wr); err != nil {
			return err
		}
	}
	return nil
}

// ReadFrom reads data from io.Reader
func (iseq *IntSeq) Read(rd io.Reader) error {
	for {
		r := newIntSeqEntry(0)
		if err := r.readFrom(rd); err != nil {
			return nil
		}
		iseq.Entries = append(iseq.Entries, *r)
	}
}

func (iseq *IntSeq) String() string {
	var s = []string{}
	for _, v := range iseq.Entries {
		s = append(s, v.String())
	}
	return "[" + strings.Join(s, ",") + "]"
}
