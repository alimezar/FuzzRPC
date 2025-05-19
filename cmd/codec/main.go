// cmd/codec/main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/alimezar/FuzzRPC/pkg/codec"
)

func main() {
	decode := flag.Bool("decode", false, "decode grpc-web-text → raw protobuf")
	encode := flag.Bool("encode", false, "encode raw protobuf → grpc-web-text")
	flag.Parse()

	if (*decode && *encode) || (!*decode && !*encode) {
		fmt.Fprintln(os.Stderr, "Choose exactly one of --decode or --encode")
		os.Exit(1)
	}

	// read stdin to EOF
	in, err := io.ReadAll(bufio.NewReader(os.Stdin))
	if err != nil {
		fmt.Fprintf(os.Stderr, "read stdin: %v\n", err)
		os.Exit(1)
	}

	var out []byte
	if *decode {
		out, err = codec.DecodeText(bytesTrimSpace(in))
	} else {
		out = codec.EncodeText(bytesTrimSpace(in))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "codec error: %v\n", err)
		os.Exit(1)
	}

	os.Stdout.Write(out)
	if !*decode { // add newline for convenience
		fmt.Println()
	}
}

// bytesTrimSpace is a tiny helper so users can pipe echo "AAA=" | codec --decode
func bytesTrimSpace(b []byte) []byte {
	for len(b) > 0 && (b[0] == '\n' || b[0] == '\r') {
		b = b[1:]
	}
	for len(b) > 0 && (b[len(b)-1] == '\n' || b[len(b)-1] == '\r') {
		b = b[:len(b)-1]
	}
	return b
}
