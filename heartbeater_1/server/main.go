package main

import (
	"time"

	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pb "heartbeater.com.server/pb"
	service "heartbeater.com.server/service"
)

func main() {

	keepaliveParams := keepalive.ServerParameters{
		Time:    1 * time.Second,
		Timeout: 5 * time.Second,
	}
	keepaliveParamsOption := grpc.KeepaliveParams(keepaliveParams)

	keepaliveEnforcementPolicy := keepalive.EnforcementPolicy{
		MinTime:             10 * time.Second,
		PermitWithoutStream: true,
	}
	keepaliveEnforcementPolicyOption := grpc.KeepaliveEnforcementPolicy(keepaliveEnforcementPolicy)

	server := grpc.NewServer(keepaliveParamsOption, keepaliveEnforcementPolicyOption)

	pb.RegisterHeartbeaterServer(server, &service.HeartbeaterService{})

	lis, _ := net.Listen("tcp", ":50051")

	server.Serve(lis)
}
