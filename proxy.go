package proxy

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"sync"
)

/*
 * This Proxy is inspired by the following:
 * https://github.com/jpillora/go-tcp-proxy
 * https://github.com/kahlys/tcpproxy
 */

// Maximum TCP size (64K)
const bufferSize = 0xffff

func prettyBytes(b uint64) string {
	var unit uint64 = 1024
	pre := [...]string{"k", "M", "G", "T", "P", "E"}
	if b < unit {
		return fmt.Sprintf("%v B", b)
	}
	exp := int(math.Log(float64(b)) / math.Log(float64(unit)))
	prefix := pre[exp-1]
	bytes := float64(b) / math.Pow(float64(unit), float64(exp))
	return fmt.Sprintf("%.1f %sB", bytes, prefix)
}

// Proxy
type Proxy struct {
	sentBytes     uint64
	receivedBytes uint64
	laddr, raddr  *net.TCPAddr
	lconn, rconn  *net.TCPConn
	mux           sync.Mutex //protect erred
	erred         bool
	done          chan bool //signal when one end of the proxy closes/errors
}

// New creates a new Proxy instance
func New(lconn *net.TCPConn, laddr, raddr *net.TCPAddr) *Proxy {
	return &Proxy{
		lconn: lconn,
		laddr: laddr,
		raddr: raddr,
		erred: false,
		done:  make(chan bool),
	}
}

/*
 * Start will open a connection to raddr
 * and wire up the connections
 */
func (p *Proxy) Start() {
	defer p.lconn.Close()

	// Setup the rconn
	var err error
	p.rconn, err = net.DialTCP("tcp", nil, p.raddr)
	if err != nil {
		//log the error
		return
	}
	defer p.rconn.Close()

	// Setup the pipes
	go p.pipe(p.lconn, p.rconn)
	go p.pipe(p.rconn, p.lconn)

	// Block until one side of the connection closes/errors
	<-p.done
	//XXX Log connection stats
	log.Printf("Connection closed (%s -> %s) sent: %s received: %s\n",
		p.rconn.LocalAddr(), p.raddr, prettyBytes(p.sentBytes),
		prettyBytes(p.receivedBytes))
}

// err is called when either lconn or rconn produce an error on read/write
func (p *Proxy) err(s string, err error) {
	// Make sure only one goroutine can cause the Proxy shutdown flow
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.erred {
		return
	}

	// EOF means connection was closed -- log anything else
	if err != io.EOF {
		log.Printf("Connection error: %s\n", err)
	}

	p.done <- true
	p.erred = true
}

// pipe writes data from the src to the dst
func (p *Proxy) pipe(src, dst *net.TCPConn) {
	// Keep track of bytes sent vs received
	isLocal := src == p.lconn

	buf := make([]byte, bufferSize)

	for {
		n, err := src.Read(buf)
		if err != nil {
			p.err("Read error '%s'\n", err)
			return
		}

		b := buf[:n]

		n, err = dst.Write(b)
		if err != nil {
			p.err("Write error '%s'\n", err)
			return
		}

		if isLocal {
			p.sentBytes += uint64(n)
		} else {
			p.receivedBytes += uint64(n)
		}
	}
}
