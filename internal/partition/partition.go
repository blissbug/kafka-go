// Package partition will wrap a storage.Log and add partition-level
// concerns: assigning a partition ID, and (once topic.go exists)
// being addressed by topic name + partition number rather than just
// a raw directory path.
//
// Not implemented yet - this is milestone 3 in the README roadmap.
package partition
