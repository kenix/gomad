# gomad
#### Golang Madness
Mike learns to Go! A journey from extreme Java to mad Golang.
* Practices with Golang
    * Prime sieves
    * Counting inversions (similarity, left, right and split inversions)
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
