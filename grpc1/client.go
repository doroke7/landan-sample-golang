package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hello/pb"
)

func main() {
	// 建立連線
	oTransportCredentials := insecure.NewCredentials()
	oTransportCredentialsOption := grpc.WithTransportCredentials(oTransportCredentials)

	oConnection, err := grpc.Dial("localhost:50051", oTransportCredentialsOption)
	if err != nil {
		log.Fatalf("無法連線: %v", err)
	}
	defer oConnection.Close()

	oClient := pb.NewGreeterClient(oConnection)

	// 呼叫服務
	oBackgroundContext := context.Background()
	ctx, cancel := context.WithTimeout(oBackgroundContext, time.Second)
	defer cancel()

	oHelloRequest := &pb.HelloRequest{Name: "Gemini"}
	oResponse, err := oClient.SayHello(ctx, oHelloRequest)
	if err != nil {
		log.Fatalf("呼叫失敗: %v", err)
	}
	sMessage := oResponse.GetMessage()
	log.Printf("伺服器回傳: %s", sMessage)
}
