package comm

import (
	"net"

	"github.com/kenix/gomad/util"
)

type Client struct {
	conn  net.Conn
	bytes uint64
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn, 0}
}

func (cln *Client) Do(dc DualChan, notify chan<- *Client) {
	for dat := range dc.Out {
		n, err := snd(cln.conn, dat)
		if err != nil {
			notify <- cln
			util.Li.Printf("notify close client from %s\n", cln.conn.RemoteAddr())
			return
		}
		cln.bytes += uint64(n)
	}
	util.Li.Printf("done on client from %s\n", cln.conn.RemoteAddr())
}

func (cln *Client) Close() {
	util.Li.Printf("close client from %s\n", cln.conn.RemoteAddr())
	cln.conn.Close()
}

func (cln *Client) BytesTransferred() uint64 {
	return cln.bytes
}

func snd(conn net.Conn, dat []byte) (int, error) {
	sent := 0
	for sent < len(dat) {
		n, err := conn.Write(dat[sent:])
		if err != nil {
			util.Le.Printf("failed sending data %s\n", err)
			return sent, err
		}
		sent += n
	}
	util.Lt.Printf("sent %d byte(s) to %s\n", sent, conn.RemoteAddr())
	return sent, nil
}
