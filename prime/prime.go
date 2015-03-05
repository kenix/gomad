package prime

type Sieve interface {
	All(uint64) <-chan uint64
}

var SieveEratosthenes Sieve = &eratosthenes{}
var SieveEratosthenesO Sieve = &eratosthenes_o{}
var SieveAtkin Sieve = &atkin{}
