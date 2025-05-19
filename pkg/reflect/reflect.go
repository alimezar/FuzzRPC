// pkg/reflect/reflect.go
package reflect

import (
	"context"

	"google.golang.org/grpc"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
)

// ListServices connects to a gRPC server via reflection and returns
// a slice of ServiceDescriptors (with their methods) or an error.
func ListServices(ctx context.Context, conn *grpc.ClientConn) ([]*desc.ServiceDescriptor, error) {
	client := grpcreflect.NewClient(
		ctx,
		reflectionpb.NewServerReflectionClient(conn),
	)
	names, err := client.ListServices()
	if err != nil {
		return nil, err
	}
	var out []*desc.ServiceDescriptor
	for _, svcName := range names {
		sd, err := client.ResolveService(svcName)
		if err != nil {
			// skip services we can’t resolve
			continue
		}
		out = append(out, sd)
	}
	return out, nil
}
