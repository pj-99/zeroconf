module github.com/pj-99/zeroconf

go 1.13

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/grandcat/zeroconf v0.0.0-00010101000000-000000000000
	github.com/miekg/dns v1.1.41
	github.com/pkg/errors v0.9.1
	golang.org/x/net v0.0.0-20210423184538-5f58ad60dda6
)

replace github.com/grandcat/zeroconf => github.com/pj-99/zeroconf v0.0.0-20250312194759-92c9e41cf3cc
