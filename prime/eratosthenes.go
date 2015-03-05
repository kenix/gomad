package prime

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
