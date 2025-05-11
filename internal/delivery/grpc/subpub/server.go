package subpub

import (
	"context"
	pb "github.com/vcreatorv/vk-sub-pub/internal/delivery/grpc/subpub/proto"
	"github.com/vcreatorv/vk-sub-pub/internal/usecase"
	"github.com/vcreatorv/vk-sub-pub/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

type GRPC struct {
	pb.UnimplementedPubSubServer
	subPub usecase.SubPub
}

func NewGRPC(uc usecase.SubPub) *GRPC {
	return &GRPC{
		subPub: uc,
	}
}

func (s *GRPC) Subscribe(r *pb.SubscribeRequest, stream grpc.ServerStreamingServer[pb.Event]) error {
	ctx := stream.Context()
	timeout := utils.GetTimeout(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sub, err := s.subPub.Subscribe(r.Key, func(msg interface{}) {
		time.Sleep(timeout)
		if err := stream.Send(&pb.Event{Data: msg.(string)}); err != nil {
			log.Printf("Stream terminated unexpectedly: %v", err)
			cancel()
		}
	})

	if err != nil {
		return status.Errorf(codes.Internal, "ошибка при подписке: %v", err)
	}

	defer sub.Unsubscribe()
	<-ctx.Done()
	return nil
}

func (s *GRPC) Publish(ctx context.Context, r *pb.PublishRequest) (*emptypb.Empty, error) {
	err := s.subPub.Publish(r.Key, r.Data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при публикации: %v", err)
	}
	return &emptypb.Empty{}, nil
}
