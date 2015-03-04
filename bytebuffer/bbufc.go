package bytebuffer

import us "unsafe"

type bbufc struct {
	bbuf
	addr uintptr
}

func NewBB(size int) ByteBuffer {
	bb := bbuf{make([]byte, size, size), 0, size, -1, ByteOrder}
	return &bbufc{bb, uintptr(us.Pointer(&bb.buf[0]))}
}

func WrapBB(bs []byte) ByteBuffer {
	bb := bbuf{bs, 0, len(bs), -1, ByteOrder}
	return &bbufc{bb, uintptr(us.Pointer(&bb.buf[0]))}
}

func (bb *bbufc) pos() us.Pointer {
	return us.Pointer(bb.addr + uintptr(bb.position))
}

const ptrSize = us.Sizeof((*byte)(nil))

func memmove(pdst, psrc us.Pointer, x int) {
	dst := uintptr(pdst)
	src := uintptr(psrc)
	n := uintptr(x)

	switch {
	case src < dst && src+n > dst:
		for i := n; i > 0; { // byte copy backward, careful: i is unsigned
			i--
			*(*byte)(us.Pointer(dst + i)) = *(*byte)(us.Pointer(src + i))
		}
	case (n|src|dst)&(ptrSize-1) != 0: // byte copy forward
		for i := uintptr(0); i < n; i++ {
			*(*byte)(us.Pointer(dst + i)) = *(*byte)(us.Pointer(src + i))
		}
	default: // word copy forward
		for i := uintptr(0); i < n; i += ptrSize {
			*(*uintptr)(us.Pointer(dst + i)) = *(*uintptr)(us.Pointer(src + i))
		}
	}
}

func (bb *bbufc) Compact() ByteBuffer {
	n := bb.Remaining()
	if n > 0 {
		memmove(us.Pointer(bb.addr), bb.pos(), n)
	}
	bb.position = n
	bb.limit = cap(bb.buf)
	bb.mark = -1
	return bb
}

func (bb *bbufc) Put(b byte) ByteBuffer {
	if bb.position >= cap(bb.buf) {
		panic(ErrOverflow)
	}
	memmove(bb.pos(), us.Pointer(&b), 1)
	bb.position++
	return bb
}

func (bb *bbufc) PutN(s []byte) ByteBuffer {
	length := len(s)
	if bb.position+length > cap(bb.buf) {
		panic(ErrOverflow)
	}
	memmove(bb.pos(), us.Pointer(&s[0]), length)
	bb.position += length
	return bb
}

func (bb *bbufc) Get() byte {
	if bb.Remaining() < 1 {
		panic(ErrUnderflow)
	}
	var b byte
	memmove(us.Pointer(&b), bb.pos(), 1)
	bb.position += 1
	return b
}

func (bb *bbufc) GetN(n int) []byte {
	if bb.Remaining() < n {
		panic(ErrUnderflow)
	}
	s := make([]byte, n)
	memmove(us.Pointer(&s[0]), bb.pos(), n)
	bb.position += n
	return s
}

func (bb *bbufc) PutUint16(i uint16) ByteBuffer {
	if bb.position+2 > cap(bb.buf) {
		panic(ErrOverflow)
	}
	ui := i
	if bb.order != ByteOrder {
		ui = i<<8 | i>>8

	}
	memmove(bb.pos(), us.Pointer(&ui), 2)
	bb.position += 2
	return bb
}

func (bb *bbufc) GetUint16() uint16 {
	if bb.Remaining() < 2 {
		panic(ErrUnderflow)
	}
	var i uint16
	memmove(us.Pointer(&i), bb.pos(), 2)
	bb.position += 2
	if bb.order != ByteOrder {
		return i<<8 | i>>8
	}
	return i
}

func (bb *bbufc) PutUint32(i uint32) ByteBuffer {
	if bb.position+4 > cap(bb.buf) {
		panic(ErrOverflow)
	}
	ui := i
	if bb.order != ByteOrder {
		ui = i<<24 | i<<8&uint32(BM16_24) | i>>8&uint32(BM08_16) | i>>24

	}
	memmove(bb.pos(), us.Pointer(&ui), 4)
	bb.position += 4
	return bb
}

func (bb *bbufc) GetUint32() uint32 {
	if bb.Remaining() < 4 {
		panic(ErrUnderflow)
	}
	var i uint32
	memmove(us.Pointer(&i), bb.pos(), 4)
	bb.position += 4
	if bb.order != ByteOrder {
		return i<<24 | i<<8&uint32(BM16_24) | i>>8&uint32(BM08_16) | i>>24
	}
	return i
}

func (bb *bbufc) PutUint64(i uint64) ByteBuffer {
	if bb.position+8 > cap(bb.buf) {
		panic(ErrOverflow)
	}
	ui := i
	if bb.order != ByteOrder {
		ui = i<<56 | i<<40&BM48_56 |
			i<<24&BM40_48 | i<<8&BM32_40 | i>>8&BM24_32 |
			i>>24&BM16_24 | i>>40&BM08_16 | i>>56
	}
	memmove(bb.pos(), us.Pointer(&ui), 8)
	bb.position += 8
	return bb
}

func (bb *bbufc) GetUint64() uint64 {
	if bb.Remaining() < 8 {
		panic(ErrUnderflow)
	}
	var i uint64
	memmove(us.Pointer(&i), bb.pos(), 8)
	bb.position += 8
	if bb.order != ByteOrder {
		return i<<56 | i<<40&BM48_56 | i<<24&BM40_48 | i<<8&BM32_40 |
			i>>8&BM24_32 | i>>24&BM16_24 | i>>40&BM08_16 | i>>56
	}
	return i
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func (bb *bbufc) Read(s []byte) (int, error) {
	n := min(bb.Remaining(), len(s))
	memmove(us.Pointer(&s[0]), bb.pos(), n)
	bb.position += n
	return n, nil
}

func (bb *bbufc) Write(s []byte) (int, error) {
	n := min(cap(bb.buf)-bb.position, len(s))
	memmove(bb.pos(), us.Pointer(&s[0]), n)
	bb.position += n
	return n, nil
}
