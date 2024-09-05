# ticrypt-file-copy [ticp]

### About
Simple CLI/Library to copy files at high speed

### Features
- High speed file copy
- Daemon mode (Accepts requests over a tcp websocket)[In testing]
- Benchmark mode
- Progress bar
- CLI interface
- Recovery Mode

### Benchmark
| Copy Mechanism | File Size | Rate of Copy            |
|----------------|-----------|-------------------------|
| DD             | 1Gb       | 220.7 MB/s              |
| DD             | 5Gb       | 224.2 MB/s              |
| Rsync          | 1GB       | 302.5 MB/s              |
| Rsync          | 5GB       | 322.1 MB/s              |
| cp             | 1GB       | 547.4 MB/s              |
| cp             | 5GB       | 558.5 MB/s              |
| ticp           | 1GB       | 919.0 MB/s              |
| ticp           | 5GB       | 1043.8 MB/s             |

### How it works
Basic principle is to use memory that is allocated using MMAP. This allows us to get page aligned memory instead of getting random memory on the heap. When we copy files using this memory we are able to levarage Direct Memory Access (DMA) which is much faster.
We are also maximising the performance by using a gorutine(for non-go users Goruoutine is a thread) for read and a separate gorotuine for write. We have a read buffer and a write buffer which are individually provisioned using MMAP. We switch between the read and write buffer when the read buffer is full and the write buffer has been written to disk. This allows us to squeeze the maximum performance out of the system when the write is slower, since we already have the read buffer ready to go by the time the write completes.

We tried to use DirectIO but as explained in this thread, we were getting much worse performance with it
https://github.com/ncw/directio/issues/2

### Configuration
The configuration file is located at `/var/lib/ticp/ticp.toml`. The configuration file is in TOML format. The configuration file has the following options:

```toml
[server]
# Host IP's to allow connections from
allowed_hosts = ["127.0.0.1"]
# Port to listen for tcp connections on
port = "4242"

[storage]
# Path to the database file
db = "ticp.db"
path = "/home/vishisht"

[copy]
# in MB
chunk_size = 4
```


### Build
```make build```

### Install
```make install```

### Usage
```
NAME:
   ticrypt-file-copy - Hight performance tool to copy files

USAGE:
   ticp [source] [destination]

COMMANDS:
   start-daemon, d  Start the daemon
   benchmark, b     Run the benchmark
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

