package service

import (
	"context"
	"time"

	pb "heartbeater.com.server/pb"
)

// HeartbeaterService 實作 pb.HeartbeaterServer 定義的 Ping 方法
type HeartbeaterService struct {
	pb.UnimplementedHeartbeaterServer
}

func (s *HeartbeaterService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	now := time.Now()
	serverTime := now.Unix()

	return &pb.PongResponse{
		Status:     "OK",
		ServerTime: serverTime,
	}, nil
}
