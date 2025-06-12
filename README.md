# TCP Traffic Server

## Overview

A simple TCP server written in Go that:
- Listens on port 9000
- Broadcasts messages to all clients except the sender
- Tracks upload/download bytes per client
- Disconnects clients after 100 bytes transferred

## How to Start

```bash
go run main.go
```

Connect using a TCP client like netcat:

```bash
nc localhost 9000
```

## How to Test

```bash
go test ./server/ -v
```