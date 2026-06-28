// Package client is the public, importable Go client for talking to
// the broker over the network protocol defined in internal/network.
// Kept separate from internal/ so it could, in principle, be `go get`-able
// by someone else's project once the wire protocol stabilizes.
//
// Not implemented yet - depends on internal/network existing first.
package client
