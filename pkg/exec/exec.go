// pkg/exec/exec.go
package exec

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alimezar/FuzzRPC/pkg/report"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mapping Severity
func classify(code codes.Code) report.Severity {
	switch code {
	case codes.Internal, codes.DataLoss:
		return report.SevCritical
	case codes.Unavailable, codes.DeadlineExceeded:
		return report.SevHigh
	case codes.InvalidArgument, codes.OutOfRange, codes.ResourceExhausted:
		return report.SevLow
	default:
		return report.SevNone
	}
}

// ExecuteFuzz invokes each mutated message against the given method,
// prints results, and calls appendFinding for each invocation.
func ExecuteFuzz(
	conn *grpc.ClientConn,
	svcName string,
	md *desc.MethodDescriptor,
	msgs []*dynamic.Message,
	appendFinding func(report.Finding),
) {
	stub := grpcdynamic.NewStub(conn)
	fullMethod := fmt.Sprintf("%s/%s", svcName, md.GetName())

	for _, msg := range msgs {
		// Prepare payload JSON for reporting
		b, err := json.Marshal(msg)
		var payload string
		if err != nil {
			payload = fmt.Sprintf("[JSON marshal error: %v]", err)
		} else {
			payload = string(b)
		}

		// Invoke RPC with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = stub.InvokeRpc(ctx, md, msg)
		cancel()

		var (
			errStr string
			sev    report.Severity
		)

		if err != nil {
			errStr = err.Error()
			sev = classify(status.Code(err))
			fmt.Printf("⚠️  %s → %v\n", fullMethod, err)
		} else {
			sev = report.SevNone
			fmt.Printf("✅  %s → OK\n", fullMethod)
		}

		// Record finding
		appendFinding(report.Finding{
			Service:   svcName,
			Method:    md.GetName(),
			Payload:   payload,
			Error:     errStr,
			Severity:  sev, // overwritten Successfully
			Baseline:  report.BaseNew,
			Timestamp: time.Now(),
		})
	}
}
