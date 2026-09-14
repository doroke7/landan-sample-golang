package bootstrap

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(dsn string) (*grpc.ClientConn, error) {
	oCredentials := insecure.NewCredentials()
	oTransportCredentials := grpc.WithTransportCredentials(oCredentials)
	oClientConn, oErr := grpc.NewClient(dsn, oTransportCredentials)

	return oClientConn, oErr
}
