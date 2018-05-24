Flux Capacitor
==============

FLoC is a set of command line tools that can be used together to implement backups.

Main features:

1. Simple: specialized replaceable processes using clean interfaces.
1. Distributed: client server architecture.
1. Performant: optimized for data safety and speed.
1. Deduplicating: data and metadata is deduplicated.
1. Incremental: backups are incrementally made but data is non-incrementally stored.

Current state
=============

Not ready for production yet.

Architecture
============

Clients read and write the files to backup and restore and send and receive data to and from the Servers.

Each Server offer the same JSON-RPC v1.0 protocol to their clients through named sockets. There is one Server per backend system that keeps the data. Planned backends:

1. `floc-leveldb`: storage resides in LevelDB.
1. `floc-boltdb`: storage resides in BoltDB.

Each Server includes a client for configuration purposes:

1. `floc-leveldb-admin`: configures a `floc-leveldb` Server.
1. `floc-boltdb-admin`: configures a `floc-boltdb` Server.

The files and metadata of a single backup are called an `Archive`. A group of `Archives` is called a `Vault`. `Vaults` are identified by a string. `Archives` belong to one `Vault` and are identified by that `Vault` and a timestamp.

A backup is performed incrementally, but the `Archive` stores the full view, effectively being a full backup. Backups can be frequent without sacrificing space because of the deduplication.

If a backup is interrupted then `Archives` are flagged as incomplete.

A stream or file of JSON documents that contain file metadata like name, type, ownership, permissions, timestamps, and path on the disk but does not include contents or extended attributes is called a `Catalog`.

The Clients that perform the actual backup and restore on any Server are:

1. `floc-read`: reads the file system and generates a `Catalog`.
1. `floc-upload`: reads a `Catalog`, partitions the contents and extended attributes of a file in chunks, deduplicates the chunks, sends only the new chunks to a Server, effectively creating an new `Archive` in the Server.
1. `floc-vault`: browses `Vaults` of a Server.
1. `floc-archive`: browses `Archives` of a Server.
1. `floc-download`: receives from a Server a `Catalog` extended with the ids of the contents and extended attributes of the files.
1. `floc-write`: reads a `Catalog`, downloads the contents of the files and writes them to the file system.
1. `floc-copy`: copies archives between servers.
1. `floc-prunable`: lists archives that may be obsolete according to some policy.
1. `floc-catalog`: reads a `Catalog` and returns a possibly different one after applying filters and transformations to the file metadata.

If a backend allows the removal of an `Archive` or a `Vault` then it must support a garbage collection mechanism to free disk storage in a way that chunks are retained only when they are 'reachable' from the remaining `Archives`. If a backend does not allow removals then a combination of `floc-prunable` and `floc-copy` may be used.

Go get them
===========

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-leveldb`

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-leveldb-admin`

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-boltdb`

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-boltdb-admin`

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-read`

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-upload`

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-vault`

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-archive`

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-download`

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-write`

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-copy`

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-prunable`

`go get github.com/daniel-fanjul-alcuten/floc/cmd/floc-catalog`
