package comm

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	bb "github.com/mentopolis/gomad/bytebuffer"
	sp "github.com/mentopolis/gomad/supplier"
	"github.com/mentopolis/gomad/util"
)

type Server struct {
	service   string
	supplier  sp.Supplier
	clients   map[*Client]DualChan
	done      chan util.Cue
	monitorCh chan *Client
	wgGroup   sync.WaitGroup
}

func NewServer(service string, supplier sp.Supplier) *Server {
	return &Server{service, supplier,
		make(map[*Client]DualChan), make(chan util.Cue),
		make(chan *Client, 1), sync.WaitGroup{}}
}

func (srv *Server) Start() error {
	util.Li.Printf("to listen @%s\n", srv.service)
	listener, err := net.Listen("tcp", srv.service)
	if err != nil {
		util.Le.Printf("failed listening @ %s: %s\n", srv.service, err)
		return err
	}
	util.Li.Printf("service ready @%s\n", srv.service)
	go srv.accept(listener)
	go srv.monitor()
	go srv.serve(listener)
	return nil
}

func (srv *Server) accept(listener net.Listener) {
	srv.wgGroup.Add(1)
	defer srv.wgGroup.Done()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "closed network") {
				util.Li.Println("stopped accepting connections")
				return
			}
			util.Le.Println(err)
		}
		util.Li.Printf("got connection from %s\n", conn.RemoteAddr())
		dc := DualChan{make(chan []byte), make(chan []byte)}
		client := NewClient(conn)
		srv.clients[client] = dc
		go client.Do(dc, srv.monitorCh)
	}
}

func (srv *Server) monitor() {
	srv.wgGroup.Add(1)
	defer srv.wgGroup.Done()
	util.Li.Println("started monitoring")
	for c := range srv.monitorCh {
		close(srv.clients[c].Out)
		c.Close()
		delete(srv.clients, c)
	}
	util.Li.Println("stopped monitoring")
}

func (srv *Server) serve(listener net.Listener) {
	buf := bb.New(1024)
	srv.wgGroup.Add(1)
	defer srv.wgGroup.Done()
	util.Li.Println("started serving")
	for {
		select {
		case <-srv.done:
			err := listener.Close()
			if err != nil {
				util.Le.Println(err)
			}
			srv.closeClients()
			close(srv.monitorCh)
			util.Li.Println("stopped serving")
			return
		case <-time.After(time.Second):
			srv.supplier.Get(buf)
			if buf.Flip().HasRemaining() {
				srv.commData(buf.GetN(buf.Remaining()))
			}
			buf.Clear()
		}
	}
}

func (srv *Server) closeClients() {
	for c := range srv.clients {
		srv.monitorCh <- c
	}
}

func (srv *Server) commData(dat []byte) {
	for _, dc := range srv.clients {
		select {
		case dc.Out <- dat:
		}
	}
}

func (srv *Server) Stop() error {
	close(srv.done)
	srv.wgGroup.Wait()
	return nil
}

func (srv *Server) Status() string {
	return fmt.Sprintf("%d client(s)", len(srv.clients))
}
