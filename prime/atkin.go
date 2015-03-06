package prime

type atkin struct{}

func (ps *atkin) All(upTo uint64) <-chan uint64 {
	r := make(chan uint64)
	go sieve_a(r, upTo)
	return r
}

func sieve_a(out chan uint64, upTo uint64) {
	// TODO
	s := []uint64{2, 3, 5, 7}
	for _, p := range s {
		out <- p
	}
	close(out)
}
