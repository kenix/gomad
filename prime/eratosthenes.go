package prime

import "math"

type eratosthenes struct{}

func (ps *eratosthenes) All(upTo uint64) <-chan uint64 {
	r := make(chan uint64)
	go sieve_e(r, upTo)
	return r
}

func generate_e(upTo uint64, ch chan uint64) {
	for i := uint64(2); i <= upTo; i++ {
		ch <- i
	}
	close(ch) // reached upper bound
}

func filter_e(p uint64, ppIn chan uint64, ppOut chan uint64) {
	for n := range ppIn {
		if n%p != 0 {
			ppOut <- n
		}
	}
	close(ppOut)
}

func sieve_e(out chan uint64, upTo uint64) {
	usch := make(chan uint64)
	go generate_e(upTo, usch)
	for {
		p, ok := <-usch
		if ok {
			out <- p
		} else {
			close(out)
			return
		}
		// TODO instead of span off new goroutine immediately, cumulate primes
		// and span off new goroutine filtering based on cumulated primes other
		// than just one
		dsch := make(chan uint64)
		go filter_e(p, usch, dsch)
		usch = dsch
	}
}

type eratosthenes_o struct{}

func (ps *eratosthenes_o) All(upTo uint64) <-chan uint64 {
	a := make([]bool, upTo-1, upTo-1) // false means prime
	rn := uint64(math.Sqrt(float64(upTo))) + 1

	r := make(chan uint64)
	go func() {
		i := uint64(2)
		for ; i < rn; i++ {
			if !a[i-2] {
				r <- i
				k := i * i
				for s, j := uint64(0), k; j <= upTo; s, j = s+1, k+s*i {
					a[j-2] = true // means composite
				}
			}
		}
		for j := i; j <= upTo; j++ {
			if !a[j-2] {
				r <- j
			}
		}
		close(r)
	}()
	return r
}
