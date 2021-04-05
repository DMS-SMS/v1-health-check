// Create file in v.1.0.0
// agent_ping.go is file that define method of gRPCAgent that agent command about ping
// For example in consul command, there are ping for connection check, etc ...

package grpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// PingToCheckConn ping for connection check to gRPC node
func (ga *gRPCAgent) PingToCheckConn(ctx context.Context, target string, opts ...grpc.DialOption) error {
	_, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
	return errors.Wrapf(err, "failed to dial gRPC with context, target: %s", target)
}
