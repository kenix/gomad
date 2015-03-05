package prime

type Sieve interface {
	All(uint64) <-chan uint64
}

var SieveEratosthenes Sieve = &eratosthenes{}
var SieveAtkin Sieve = &atkin{}
