// package prime implements several prime number sieves
package prime

// Sieve gets an upper bound as input and returns a channel, from which all prime
// numbers upto that upper bound can be drained, once.
type Sieve interface {
	All(uint64) <-chan uint64
}

var sieveEratosthenes Sieve = &eratosthenes{}

// SieveEratosthenes implements the sieve of Eratosthenes using exclusion-channel
// algorithm adapted from Rob Pike's presentation of Concurrency Pattern
func SieveEratosthenes() Sieve {
	return sieveEratosthenes
}

var sieveEratosthenesO Sieve = &eratosthenes_o{}

// SieveEratosthenesO implements the optimized Eratosthenes sieve
func SieveEratosthenesO() Sieve {
	return sieveEratosthenesO
}

var sieveAtkin Sieve = &atkin{}

// TODO not implemented yet
func SieveAtkin() Sieve {
	return sieveAtkin
}

var sieveSundaram Sieve = &sundaram{}

// SieveSundaram implements the prime sieve of Sundaram
func SieveSundaram() Sieve {
	return sieveSundaram
}
