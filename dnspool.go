// Package dnspool uses a goroutine pool for DNS resolving to limit number of
// OS threads that will be spawned by the go runtime.
//
// It will nolonger be needed when the go runtime provides a way to limit OS
// threads creation.

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

var requestChan chan *lookupRequest

// Initialize the DNS resove goroutine pool with n goroutines. Must call this
// first using any other functions in this package.
func InitDNSPool(n int) {
	const defaultNGoroutine = 100

	if n <= 0 {
		fmt.Printf("initDNSPool parameter error: %d is not positive, using default value %d\n",
			n, defaultNGoroutine)
	} else {
		n = defaultNGoroutine
	}

	requestChan = make(chan *lookupRequest, n)
	for i := 0; i < n; i++ {
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

// Same as net.LookupHost. Not thread safe.
func (r *Resolver) LookupHost(host string) (addrs []string, err error) {
	r.host = host
	requestChan <- (*lookupRequest)(r)
	<-r.done
	return r.addrs, r.err
}

// Same as net.LookupHost. If you are going to call this repeatedly in the
// same goroutine, it's better to create a new Resolver to avoid some
// performance overhead.
func LookupHost(host string) (addrs []string, err error) {
	r := NewResolver()
	return r.LookupHost(host)
}
