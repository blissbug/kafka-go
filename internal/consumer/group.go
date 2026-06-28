// Package consumer will implement consumer groups: coordinating
// multiple consumers so each partition is only being read by one
// member of the group at a time, plus tracking committed offsets so
// a consumer can resume after restarting instead of re-reading
// everything from the start.
//
// This is milestone 4 in the README roadmap, and the one flagged as
// the area to go deep on for differentiation - most "Kafka from
// scratch" projects skip real rebalancing logic, so doing this one
// properly (handling a consumer joining/leaving the group mid-stream)
// is worth the extra time.
package consumer
