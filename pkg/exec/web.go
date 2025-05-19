package exec

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/alimezar/FuzzRPC/pkg/codec"
	"google.golang.org/protobuf/proto"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// callWebUnary encodes req (proto bytes) as grpc-web-text, POSTs it, and returns proto response bytes.
func callWebUnary(ctx context.Context, endpoint, fullMethod string, pbReq []byte) ([]byte, error) {
	// POST <endpoint>/<fullMethod>   e.g. http://host:8080/grpc.gateway.testing.EchoService/Echo
	url := endpoint + "/" + fullMethod
	body := codec.EncodeText(pbReq)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/grpc-web-text")
	req.Header.Set("x-grpc-web", "1")    // some servers require it
	req.Header.Set("grpc-timeout", "5S") // optional
	req.Header.Set("Accept", "application/grpc-web-text")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return codec.DecodeText(rawResp)
}

// helper that marshals any proto.Message
func marshalProto(msg proto.Message) ([]byte, error) { return proto.Marshal(msg) }
