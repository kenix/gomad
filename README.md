# gomad
#### Golang Madness

[![Build Status](https://travis-ci.org/kenix/gomad.svg)](https://travis-ci.org/kenix/gomad)

Mike learns to Go! A journey from Java to Golang.
* Practices with Golang
    * Prime sieves
    * Algorithms refresh
    	* Counting inversions (similarity, left, right and split inversions)
    	* Closest pairs
    	* randomized quick sort util
    * concurrency patterns `Runner`, `Pool` and `Worker` from  Go in Action
* A take on `ByteBuffer/StringBuilder` and `slice/bytes.Buffer`
	* fix-sized `ByteBuffer` vs. dynamically sized `slice`
	* panic with error vs. multiple return values including error
	* method chaining vs. `append` and `copy` or util method with interface{}
	* object.oriented vs. functional
	* setter methods vs. `xxxTo`
	* honor the interfaces in common packages like `io`
    * impl. based on `slice`: w/ or w/o re-slicing benchmark comparison: how to make it faster? Too much overhead with cgo
* TCP communication in Golang
	* infinite data server (financial ticks)
	* data client
* Simple key-block database - sdb
