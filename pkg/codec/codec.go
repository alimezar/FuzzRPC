package codec

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

// --- Public API -----------------------------------------------------------

// DecodeText converts a grpc-web-text payload → raw protobuf bytes (single frame).
func DecodeText(in []byte) ([]byte, error) {
	// 1) base64-decode
	raw := make([]byte, base64.StdEncoding.DecodedLen(len(in)))
	n, err := base64.StdEncoding.Decode(raw, in)
	if err != nil {
		return nil, fmt.Errorf("base64: %w", err)
	}
	raw = raw[:n]

	// 2) validate framing: 1-byte flags | 4-byte len | message
	if len(raw) < 5 {
		return nil, fmt.Errorf("frame too short")
	}
	msgLen := binary.BigEndian.Uint32(raw[1:5])
	if int(msgLen) != len(raw[5:]) {
		return nil, fmt.Errorf("length mismatch (%d vs %d)", msgLen, len(raw[5:]))
	}
	return raw[5:], nil
}

// EncodeText wraps protobuf bytes into grpc-web frame then base64-encodes.
func EncodeText(pb []byte) []byte {
	var buf bytes.Buffer
	// 0x00 flags (no compression) + 4-byte len
	buf.WriteByte(0x00)
	binary.Write(&buf, binary.BigEndian, uint32(len(pb)))
	buf.Write(pb)
	out := make([]byte, base64.StdEncoding.EncodedLen(buf.Len()))
	base64.StdEncoding.Encode(out, buf.Bytes())
	return out
}
