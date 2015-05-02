package sdb

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/user"
	"sort"
	"strings"
	"testing"
)

func tmpDir() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return u.HomeDir + "/tmp/sdb"
}

func tmpFile(t *testing.T) string {
	f, err := ioutil.TempFile(tmpDir(), "sdb-")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	return f.Name()
}

var maxMockDatLen = int32(1 << 12) //4K
//var maxMockDatLen = int32(8)
var entryCount = 1 << 12 // 4K
//var entryCount=10

func TestAccess(t *testing.T) {
	keys := mockKeys(entryCount)
	dat := make([]*datMock, 0, entryCount)
	for _, k := range keys {
		dat = append(dat, &datMock{k, rand.Int31n(maxMockDatLen) + 1, mockDat()})
	}

	fn := tmpFile(t)
	writeDat(fn, dat, t)
	defer os.Remove(fn)
	r, err := NewReader(fn)
	if err != nil {
		t.Fatal(err)
	}

	for _, d := range dat {
		bs, err := r.Get(d.key)
		if err != nil {
			t.Errorf("failed reading data for %s\n", d.key)
		}
		if len(bs) != int(d.length) {
			t.Errorf("data len [%s] wanted %d, got %d\n", d.key, d.length, len(bs))
			continue
		}
		for i := 0; i < int(d.length); i++ {
			if bs[i] != d.dat {
				t.Errorf("data wanted %s, got %s\n", string(d.dat), string(bs[i]))
			}
		}
	}

	defer r.Close()
}

func writeDat(fn string, dat []*datMock, t *testing.T) {
	w, err := NewWriter(fn)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	for _, d := range dat {
		n, err := w.Put(d.key, d.data())
		if err != nil {
			t.Fatal(err)
		}
		if int32(n) != d.length {
			t.Fatalf("partial write %d expected, %d written", d.length, n)
		}
	}
}

var alpha = []byte("abcdefghijklmnopqrstuvwxyz")
var alphA = []byte(strings.ToUpper(string(alpha)))

func mockDat() byte {
	x := rand.Intn(26)
	b := rand.Intn(2)
	if b&0x01 == 0 {
		return alphA[x]
	} else {
		return alpha[x]
	}
}

type datMock struct {
	key    string
	length int32
	dat    byte
}

func (dm *datMock) String() string {
	return fmt.Sprintf("%s;%d;%s", dm.key, dm.length, string(dm.dat))
}

func (dm *datMock) data() (d []byte) {
	d = make([]byte, dm.length, dm.length)
	for i, _ := range d {
		d[i] = dm.dat
	}
	return
}

func mockKey() string {
	n := rand.Intn(32) + 1
	ba := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		x := rand.Intn(26)
		b := rand.Intn(2)
		if b&0x01 == 0 {
			ba = append(ba, alphA[x])
		} else {
			ba = append(ba, alpha[x])
		}
	}
	return string(ba)
}

func mockKeys(size int) []string {
	keys := make([]string, 0, 0)
	for len(keys) < size {
		k := mockKey()
		if len(keys) > 0 {
			idx := sort.SearchStrings(keys, k)
			if idx < len(keys) && keys[idx] == k {
				continue
			}
		}
		keys = append(keys, k)
		sort.Strings(keys)
	}
	return keys
}
