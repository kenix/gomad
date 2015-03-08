package comm

type DualChan struct {
	In  chan []byte
	Out chan []byte
}
