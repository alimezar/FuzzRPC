# FuzzRPC

**A native Go fuzzer for gRPC and gRPC-Web services**  
FuzzRPC discovers services via reflection, builds seed messages automatically, mutates them type-safely, and reports findings in structured JSON and a color-coded HTML dashboard.

---

## Key Features

| Area              | Capability |
|-------------------|------------|
| **Discovery**     | Reflection-driven service and method enumeration (no proto files needed) |
| **Fuzzing**       | Seed generation → type-aware field mutation → concurrent execution |
| **Transports**    | • HTTP/2 gRPC (`application/grpc`)<br>• HTTP/1.1 gRPC-Web-Text (`application/grpc-web-text`) |
| **Reporting**     | • `out.json` machine-readable log<br>• `out.html` interactive dashboard with severity tint, baseline diff, and Chart.js bar graph |
| **Diffing**       | `--baseline` flag highlights new, unchanged, and resolved findings between scans |
| **Severity**      | Maps gRPC status codes to `critical / high / low / none` |
| **CLI Helpers**   | `cmd/codec` encodes/decodes gRPC-Web-Text frames for manual testing or Burp Suite integration |
| **Zero Dependencies** | Ships as a single static Go binary (`go install ...`) |

---

## Installation

```bash
go install github.com/alimezar/FuzzRPC/cmd/fuzzrpc@latest
go install github.com/alimezar/FuzzRPC/cmd/codec@latest   # optional helper
```

> Requires Go 1.22 or newer.

---

## Quick Start

### 1. Run a target gRPC server

```bash
go run examples/helloworld/server/main.go   # listens on :50051
```

### 2. Native gRPC fuzzing

```bash
fuzzrpc \
  --target localhost:50051 \
  --report-json out.json \
  --report-html out.html
```

### 3. gRPC-Web fuzzing (via proxy on :8080)

```bash
fuzzrpc \
  --target localhost:8080 \
  --web \
  --report-json out_web.json
```

### 4. Baseline diff (CI regression gate)

```bash
fuzzrpc \
  --target staging.internal:50051 \
  --baseline previous.json \
  --report-json current.json \
  --fail-on new,critical     # forthcoming flag
```

---

## Command-line Flags

| Flag                | Description                                     | Default |
|---------------------|-------------------------------------------------|---------|
| `--target`          | host:port of the gRPC or gRPC-Web endpoint      | —       |
| `--timeout`         | Dial/call timeout                               | 5s      |
| `--web`             | Use gRPC-Web-Text transport                     | false   |
| `--report-json`     | Path to write `out.json`                        | —       |
| `--report-html`     | Path to write `out.html`                        | —       |
| `--report-template` | Custom HTML template path                       | `templates/report.html` |
| `--baseline`        | Previous `out.json` file for diffing            | —       |

Run `fuzzrpc -h` to view all available options.

---

## HTML Dashboard

- Rows are tinted by severity:
  - **Critical** → Red
  - **High** → Orange
  - **Low** → Green
  - **None** → Plain

- Left border indicates baseline status:
  - **Blue** = New
  - **Grey** = Unchanged
  - **Strike-through** = Resolved

- Includes a Chart.js bar graph summarizing findings by severity.

---

## Codec Utility

```bash
# Encode raw protobuf → gRPC-Web-Text
cat request.bin | codec --encode > payload.txt

# Decode intercepted payload
cat payload.txt | codec --decode > request.bin
```

Useful for Burp Suite or manual replay.

---

## Architecture Overview

```text
reflection → seed → mutate → (gRPC | gRPC-Web) runner → findings → report
               ↑            concurrent goroutines                ↑
               └────────────── baseline diff & severity mapping ─┘
```

Each stage lives in its own `pkg/` sub-module:
- `reflect`
- `seed`
- `mutate`
- `exec`
- `codec`
- `report`

---

## Roadmap

- Add `application/grpc-web+proto` (binary) support
- Migrate CLI to Cobra (`enum / seed / fuzz / report` subcommands)
- Burp extension (zero-dependency tab leveraging `codec`)
- Plugin system (Go‐plugin mutators and authentication hooks – JWT, mTLS)
- GitHub Actions + full test coverage
- Interactive TUI with step-through fuzzing and payload inspector

---

## Contributing

1. Fork and create a feature branch.
2. Ensure `go test ./...` passes.
3. Submit a PR with a concise description.

---

## License

Distributed under the MIT License. See [`LICENSE`](./LICENSE) for details.
