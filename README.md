# go-kafka

A Kafka clone built from scratch in Go, as a learning project. The goal
isn't to replace Kafka - it's to actually understand *why* Kafka is
built the way it is, by building each piece myself: the append-only
storage engine, the network protocol, partitioning, consumer groups,
and (eventually) replication.

## Status

**This is a work in progress, built incrementally and in public.**
Below is an honest list of what's actually implemented vs. what's
still ahead, updated as the project moves forward.

### Done
- [x] Folder structure and project layout
- [x] Append-only log storage (`internal/storage`)

### Not yet built
- [ ] Segment files with size-based rollover
- [ ] Offset index per segment for O(1) lookup by offset (no scanning)
- [ ] Reopening an existing log directory resumes correctly after a restart
- [ ] Network layer (TCP server, wire protocol) - `internal/network` exists
      as an empty package, nothing implemented yet
- [ ] Partitioning / topic abstraction - `internal/partition`, `internal/topic`
      are empty packages
- [ ] Broker tying it all together - `internal/broker` is empty
- [ ] Consumer groups + offset commit tracking - `internal/consumer` is empty
- [ ] Replication (leader/follower) - not started
- [ ] Benchmarks comparing naive vs. batched writes - not started

### Known gaps / honesty notes
- Crash recovery is naive: if the process dies mid-write to the final
  segment, behavior on reopen is untested and likely incorrect. Real
  Kafka has to handle a truncated/corrupt tail segment carefully -
  this is a good thing to tackle deliberately later, not pretend is
  solved.
- `Log.Append` holds a single mutex for the whole call, so writes are
  fully serialized. This is the deliberately naive v1 baseline -
  improving this (and measuring the improvement) is planned as part
  of the differentiation work mentioned below.
- `segmentFor` is a linear scan over all segments. Fine at small scale,
  a known place to optimize (binary search) once there's a reason to.

## Project structure

```
go-kafka/
├── cmd/server/          # entrypoint (not yet implemented)
├── internal/
│   ├── storage/          # append-only log, segments, index - IMPLEMENTED
│   ├── partition/        # partition abstraction - not yet implemented
│   ├── topic/             # topic + key-based partitioning - not yet implemented
│   ├── broker/             # ties storage+network+partition together - not yet implemented
│   ├── network/             # TCP server, wire protocol - not yet implemented
│   └── consumer/             # consumer groups, offset tracking - not yet implemented
└── pkg/client/                # public Go client - not yet implemented
```

## Why these design choices
1. Used 4 bytes for the offset index, Because a single log segment rarely holds more than 4.29 billion messages, the resulting relative offset safely fits within 4 bytes.
2. +----------------+----------------+
   | relativeOffset | bytePosition   |
   |    uint32      |    uint64      |
   +----------------+----------------+

## Roadmap

1. Storage engine 
2. Network layer + basic produce/consume over TCP
3. Topics with multiple partitions, key-based hashing
4. Consumer groups + offset commits (the part most "Kafka from scratch"
   projects skip or do shallowly - planned area of depth for this project)
5. Replication (leader/follower)

Each milestone will get its own benchmark where relevant (e.g. naive
single-write-per-message vs. batched writes, once the network layer
exists to actually drive load).
