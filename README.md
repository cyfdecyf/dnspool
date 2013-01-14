Package `dnspool` creates a goroutine pool for DNS query to limit the number of
OS threads that will be spawned by the go runtime.

Default number of goroutines for DNS lookup is set to 32 now. You can increase goroutine number by calling `dnspool.SetGoroutineNumber`. Note the maximum number of goroutines is 256. (Please tell me if this number is not enough.)

This package will no longer be needed when the go runtime provides a way to limit OS threads creation.

`dnspool.LookupHost` has the same usage as `net.LookupPost`. For example:

    addrs, err := dnspool.LookupHost("google.com")

If you are going to call `dnspool.LookupHost` many times in the same goroutine, it's better to create a `dnspool.Resolver` and call the `LookupHost` method to avoid some performance overhead.

    resolver := dnspool.NewResolver()
    addrs, err := resolver.LookupHost(host)