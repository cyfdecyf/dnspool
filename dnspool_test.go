package dnspool

import (
	"testing"
)

func init() {
	// Test SetGoroutineNumber
	SetGoroutineNumber(1)
	SetGoroutineNumber(1024)
	SetGoroutineNumber(64)
}

func checkLocalhostIP(addrs []string, err error, msg string, t *testing.T) {
	if err != nil {
		t.Error(msg, "lookup localhost error:", err)
	}
	found := false
	for _, ad := range addrs {
		if ad == "127.0.0.1" {
			found = true
		}
	}
	if !found {
		t.Error(msg, "localhost didn't resolve to 127.0.0.1")
	}
}

func TestLookupHost(t *testing.T) {
	addrs, err := LookupHost("localhost")
	checkLocalhostIP(addrs, err, "TestLookupHost", t)
}

func TestConcurrentLookupHost(t *testing.T) {
	const nResolver = 20
	const nLookup = 100

	resolve := func() {
		r := NewResolver()
		for i := 0; i < nLookup; i++ {
			addrs, err := r.LookupHost("127.0.0.1")
			checkLocalhostIP(addrs, err, "TestConcurrentLookupHost", t)
		}
	}
	for i := 0; i < nResolver; i++ {
		go resolve()
	}
}
