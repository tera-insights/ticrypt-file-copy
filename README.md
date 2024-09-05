# ticrypt-file-copy

### About
Simple CLI/Library to copy files at high speed

### Features
- High speed file copy
- Daemon mode
- Benchmark mode
- Progress bar
- CLI interface
- Recovery Mode

### Benchmark
| Copy Mechanism | File Size | Rate of Copy            |
|----------------|-----------|-------------------------|
| DD             | 1Gb       | 220.7 MB/s              |
| Rsync          | 1GB       | 302.5 MB/s              |
| cp             | 1GB       | 547.4 MB/s              |
| ticp           | 1GB       | 919.0 MB/s              |



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

