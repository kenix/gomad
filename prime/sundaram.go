package prime

type sundaram struct{}

func (ps *sundaram) All(upTo uint64) <-chan uint64 {
	r := make(chan uint64)
	go sieve_s(r, upTo)
	return r
}

func sieve_s(out chan uint64, upTo uint64) {
	n := upTo>>1 - 1
	s := make([]bool, n, n)
	for i := uint64(1); i <= n; i++ {
		for j := uint64(1); j <= i; j++ {
			k := i + j + i*j<<1
			if k <= n {
				s[k-1] = true
			}
		}
	}
	out <- uint64(2)
	for i, v := range s {
		if !v {
			out <- uint64((i+1)<<1 + 1)
		}
	}

	close(out)
}
