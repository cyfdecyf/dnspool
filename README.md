Package dnspool uses a goroutine pool for DNS resolving to limit number of
OS threads that will be spawned by the go runtime.

It will nolonger be needed when the go runtime provides a way to limit OS
threads creation.
