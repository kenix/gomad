package sdb

import (
	"os"
	"sort"

	bb "github.com/kenix/gomad/bytebuffer"
)

type rImpl struct {
	f       *os.File
	offsets []int64
	keys    []string
	curPage int32
	buf     bb.ByteBuffer
}

func NewReader(fn string) (Reader, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	// read indices start
	indicesStartPos, err := f.Seek(-8, 2)
	if err != nil {
		return nil, err
	}

	buf := bb.New(8)
	n, err := buf.ReadFrom(f)
	if err != nil || n != 8 {
		return nil, err
	}

	indicesStart := int64(buf.Flip().GetUint64())

	// read indices
	buf = bb.New(int(indicesStartPos - indicesStart))
	if _, err = f.Seek(indicesStart, 0); err != nil {
		return nil, err
	}

	for buf.HasRemaining() { // read indices fully into buffer
		if _, err = buf.ReadFrom(f); err != nil {
			return nil, err
		}
	}

	// parse offsets and keys in indices
	offsets := make([]int64, 0, 0)
	keys := make([]string, 0, 0)
	buf.Flip()
	for buf.HasRemaining() {
		offsets = append(offsets, int64(buf.GetUint64()))
		if buf.HasRemaining() {
			kl := buf.Get()
			keys = append(keys, string(buf.GetN(int(kl))))
		}
	}
	offsets = append(offsets, indicesStart) // guard offset
	return &rImpl{f, offsets, keys, -1, bb.New(size_page_key)}, nil
}

func (r *rImpl) Underlying() string {
	return r.f.Name()
}

func (r *rImpl) Close() error {
	return r.f.Close()
}

func (r *rImpl) Get(key string) ([]byte, error) {
	idx := 0
	if len(r.keys) > 0 {
		idx = sort.SearchStrings(r.keys, key)
		if idx < len(r.keys) && r.keys[idx] == key {
			idx += 1
		}
	}

	if r.curPage != int32(idx) { // read keys and OaL into page buffer
		offset := r.offsets[idx]
		length := r.offsets[idx+1] - offset
		_, err := r.f.Seek(offset, 0)
		if err != nil {
			return nil, err
		}
		r.buf.Clear()
		for r.buf.Position() < int(length) {
			if _, err := r.buf.ReadFrom(r.f); err != nil {
				return nil, err
			}
		}
		r.buf.LimitTo(int(length)).PositionTo(0)
		r.curPage = int32(idx)
	}

	// search page buffer for key
	defer r.buf.PositionTo(0)
	for r.buf.HasRemaining() {
		kl := int(r.buf.Get())
		candidate := string(r.buf.GetN(kl))
		oal := r.buf.GetUint64()
		if key == candidate {
			offset := int64(oal >> dataLengthBits)
			length := int32(oal & dataLengthMask)
			return r.loadDat(offset, length)
		}
	}

	return nil, nil
}

func (r *rImpl) loadDat(offset int64, length int32) ([]byte, error) {
	if _, err := r.f.Seek(offset, 0); err != nil {
		return nil, err
	}
	dat := make([]byte, length, length)
	for n := 0; n < len(dat); {
		i, err := r.f.Read(dat[n:])
		if err != nil {
			return nil, err
		}
		n += i
	}
	return dat, nil
}
