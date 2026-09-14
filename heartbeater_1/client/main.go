package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff" // 關鍵：引入 gRPC 官方的退避參數庫
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	pb "heartbeater.com.client/pb"
)

func main() {
	// 為了能快速看到「風暴」，我們把時間縮短：2秒送一次心跳，1秒沒回就判定超時
	transportCredentials := insecure.NewCredentials()
	transportCredentialsOption := grpc.WithTransportCredentials(transportCredentials)

	connectParams := grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  1.0 * time.Second, // 第一次斷線後，等 1.0 秒再嘗試重連
			Multiplier: 1.6,               // 每次重連失敗，等待時間乘以 1.6 (1s -> 1.6s -> 2.56s)
			Jitter:     0.2,               // 加上 20% 的隨機抖動誤差，把大量 Client 的重連時間錯開
			MaxDelay:   10 * time.Second,  // 不管失敗幾次，最長只等 30 秒，避免時間被無限拉長
		},
		MinConnectTimeout: 3 * time.Second, // 每次嘗試建立 TCP 握手時，最少給底層 3 秒的超時時間
	}
	connectParamsOption := grpc.WithConnectParams(connectParams)

	keepaliveClientParams := keepalive.ClientParameters{
		Time:                10 * time.Second, // 每 10 秒偷偷送一次 PING 保活、防止被防火牆剪斷
		Timeout:             3 * time.Second,  // PING 出去後 3 秒內 Server 沒回應，直接判定斷線，立刻觸發上面的 Backoff 流程
		PermitWithoutStream: true,             // 關鍵：就算現在業務沒請求、沒有 Stream，也要送 PING
	}
	keepaliveParamsOption := grpc.WithKeepaliveParams(keepaliveClientParams)

	/*
		重連必須要要 client 端主動發起 業務連接的時候的瞬間嗎，如果業務都沒有新請求就會持續idle


	*/
	conn, err := grpc.NewClient(
		"127.0.0.1:50051",
		transportCredentialsOption,
		connectParamsOption,
		keepaliveParamsOption,
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 關鍵動作 1：強制讓 Client 去連線，把狀態從 IDLE 推向 READY
	conn.Connect()

	client := pb.NewHeartbeaterClient(conn)

	// 關鍵動作 2：應用層心跳，2 秒送一次 Ping，1 秒沒回就判定超時
	go func() {
		for {
			backgroundContext := context.Background()
			ctx, cancel := context.WithTimeout(backgroundContext, 1*time.Second)

			now := time.Now()
			pingRequest := &pb.PingRequest{
				ClientId:  "client-1",
				Timestamp: now.Unix(),
			}
			resp, err := client.Ping(ctx, pingRequest)
			cancel()

			if err != nil {
				failureTime := time.Now()
				failureTimeText := failureTime.Format("15:04:05")
				fmt.Printf("[%s] 心跳失敗: %v\n", failureTimeText, err)
			} else {
				successTime := time.Now()
				successTimeText := successTime.Format("15:04:05")
				fmt.Printf("[%s] 心跳成功: status=%s server_time=%d\n", successTimeText, resp.Status, resp.ServerTime)
			}

			time.Sleep(2 * time.Second)
		}
	}()

	// 關鍵動作 3：必須起一個活動的 Goroutine 監控，讓 Go 知道這不是死鎖
	go func() {
		for {
			state := conn.GetState()
			checkTime := time.Now()
			checkTimeText := checkTime.Format("15:04:05")
			fmt.Printf("[%s] 當前連線狀態: %s\n", checkTimeText, state)

			// 如果不小心掉回 IDLE，就再逼它連線
			if state == connectivity.Idle {
				// conn.Connect()
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// 讓程式常駐
	select {}
}
