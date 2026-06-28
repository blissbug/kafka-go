// Package network implements the TCP server and wire protocol that
// clients use to produce and consume records. This is milestone 2 in
// the README roadmap - the next thing to build after storage.
//
// Design questions to resolve when building this:
//   - Wire format: length-prefixed framing is the simplest correct
//     choice
//   - One connection per client vs. a connection pool - start with
//     one-goroutine-per-connection, the standard simple Go starting
//     point, and only complicate this if a real bottleneck shows up
//     in a benchmark.
package network
