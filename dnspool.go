/*
Package dnspool creates a goroutine pool for DNS query to limit the number of
OS threads that will be spawned by the go runtime.

Default number of goroutines for DNS lookup is set to 32 now. You can increase
goroutine number by calling SetGoroutineNumber. Note the maximum
number of goroutines is 256. (Please tell me if this number is not enough.)

This package will no longer be needed when the go runtime provides a way to
limit OS threads creation.
*/
package dnspool

import (
	"fmt"
	"net"
)

type lookupRequest struct {
	host  string
	addrs []string // lookup result
	err   error
	done  chan byte // channel notify that lookup is done
}

// Resolver provides LookupHost method which can be used non-concurrently.
type Resolver lookupRequest

// Maximum concurrent lookup request
const maxLookupReq = 256

var requestChan = make(chan *lookupRequest, maxLookupReq)

// Default goroutine pool size. Making this larger than maxLookupReq is
// useless as the channel can hold only that many lookup request.
var nGoroutine = 32

func init() {
	for i := 0; i < nGoroutine; i++ {
		go lookup()
	}
}

// SetGoroutineNumber sets the number of goroutines used to do DNS query.
// If n <= current goroutine number or is larger than maximum concurrent DNS
// request, this function will do nothing.
func SetGoroutineNumber(n int) {
	if n <= nGoroutine {
		fmt.Printf("SetGoroutineNumber: %d <= current goroutine number %d, do nothing\n", n, nGoroutine)
		return
	}
	if n > maxLookupReq {
		fmt.Printf("SetGoroutineNumber: %d > maximum goroutines number %d, do nothing\n", n, maxLookupReq)
		return
	}
	for ; nGoroutine < n; nGoroutine++ {
		go lookup()
	}
}

func lookup() {
	for {
		req := <-requestChan
		req.addrs, req.err = net.LookupHost(req.host)
		req.done <- 1
	}
}

func NewResolver() *Resolver {
	return &Resolver{done: make(chan byte)}
}

// LookupHost has the same usage as net.LookupHost. Not thread safe.
func (r *Resolver) LookupHost(host string) (addrs []string, err error) {
	r.host = host
	requestChan <- (*lookupRequest)(r)
	<-r.done
	return r.addrs, r.err
}

// LookupHost has the same usage as net.LookupHost. If you are going to call
// this repeatedly in the same goroutine, it's better to create a new Resolver
// to avoid some performance overhead.
func LookupHost(host string) (addrs []string, err error) {
	return NewResolver().LookupHost(host)
}

// Dial has the same usage as as net.Dial. This function will use LookupHost
// to resolve host address, then try net.Dial on each returned ip address till
// one succeeds or all fail.
func Dial(hostPort string) (c net.Conn, err error) {
	var addrs []string
	var host, port string

	if host, port, err = net.SplitHostPort(hostPort); err != nil {
		return
	}
	// No need to call LookupHost if host part is IP address
	if ip := net.ParseIP(host); ip != nil {
		return net.Dial("tcp", hostPort)
	}
	if addrs, err = LookupHost(host); err != nil {
		return
	}
	for _, ip := range addrs {
		ipHost := net.JoinHostPort(ip, port)
		if c, err = net.Dial("tcp", ipHost); err == nil {
			return
		}
	}
	return nil, err
}
