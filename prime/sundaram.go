package prime

type sundaram struct{}

func (ps *sundaram) All(upTo uint64) <-chan uint64 {
	r := make(chan uint64)
	go sieve_s(r, upTo)
	return r
}

func sieve_s(out chan uint64, upTo uint64) {
	// TODO
	s := []uint64{2, 3, 5, 7}
	for _, p := range s {
		out <- p
	}
	close(out)
}
