package interceptors

import (
	"context"
	"github.com/vcreatorv/vk-sub-pub/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

func TimeoutClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		timeout := utils.GetTimeout(ctx)
		ctx = metadata.AppendToOutgoingContext(ctx, "timeout", timeout.String())
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func TimeoutServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if ok {
			if values := md.Get("timeout"); len(values) > 0 {
				timeout, err := time.ParseDuration(values[0])
				if err != nil {
					return err
				}
				newCtx := utils.SetTimeout(ss.Context(), timeout)
				wrapped := &wrappedServerStream{
					ServerStream: ss,
					ctx:          newCtx,
				}
				return handler(srv, wrapped)
			}
		}
		return handler(srv, ss)
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
