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
		dsch := make(chan uint64)
		go filter_e(p, usch, dsch)
		usch = dsch
	}
}

type eratosthenes_o struct{}

func (ps *eratosthenes_o) All(upTo uint64) <-chan uint64 {
	a := make([]bool, upTo, upTo)
	a[0], a[1] = true, true
	rn := uint64(math.Sqrt(float64(upTo))) + 1

	for i := uint64(2); i < rn; i++ {
		if !a[i] { // does not have factor other than 1 and itself
			k := i * i
			for s, j := uint64(0), k; j < upTo; s, j = s+1, k+s*i {
				a[j] = true
			}
		}
	}

	r := make(chan uint64)
	go func() {
		for i, b := range a {
			if !b {
				r <- uint64(i)
			}
		}
		close(r)
	}()
	return r
}
