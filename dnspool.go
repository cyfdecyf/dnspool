// Package dnspool creates a goroutine pool for DNS resolving to limit the
// number of OS threads that will be spawned by the go runtime.
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

// Set the number of goroutines used to do DNS query. If n <= current
// goroutine number or is larger than maximum concurrent DNS request, this
// function will do nothing.
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
	return NewResolver().LookupHost(host)
}
