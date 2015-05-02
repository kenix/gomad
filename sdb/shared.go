package sdb

import "fmt"

type entry struct {
	key    string
	offset int64 // 40 bits used, max 1T
	length int32 // 24 bits used, max 16M, offset, length together 8 bytes
}

func (e *entry) String() string {
	return fmt.Sprintf("%s;%d;%d", e.key, e.offset, e.length)
}

type entries []*entry

func (es entries) Len() int {
	return len(es)
}

func (es entries) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es entries) Less(i, j int) bool {
	return es[i].key < es[j].key
}

func (es entries) totalLength() int64 {
	n := int64(0)
	for _, e := range es {
		n += int64(e.length)
	}
	return n
}
