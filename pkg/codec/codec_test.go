package codec

import "testing"

func TestRoundTrip(t *testing.T) {
	in := []byte{0x08, 0x96, 0x01} // HelloRequest{id:150}
	enc := EncodeText(in)
	dec, err := DecodeText(enc)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if string(dec) != string(in) {
		t.Fatalf("round-trip mismatch: %#v vs %#v", dec, in)
	}
}
