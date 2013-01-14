Package dnspool uses a goroutine pool for DNS resolving to limit number of
OS threads that will be spawned by the go runtime.

It will nolonger be needed when the go runtime provides a way to limit OS
threads creation.

`dnspool.LookupHost` has the same usage as `net.LookupPost`. For example:

    addrs, err := dnspool.LookupHost("google.com")

If you are going to call `dnspool.LookupHost` many times in the same goroutine, it's better to create a `dnspool.Resolver` and call the `LookupHost` method to avoid some performance overhead.

    resolver := dnspool.NewResolver()
    addrs, err := resolver.LookupHost(host)