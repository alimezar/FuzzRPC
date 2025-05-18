# FuzzRPC

**FuzzRPC** is a CLI tool for security engineers to automatically enumerate and fuzz gRPC services via reflection.

## Features
- Reflection-driven service/method discovery  
- Structured “seed” message generation from protobuf descriptors  
- Type-aware mutation/fuzz strategies  
- Concurrent execution over HTTP/2  
- JSON and HTML reporting  

## Installation

```bash
go install github.com/alimezar/fuzzrpc@latest
