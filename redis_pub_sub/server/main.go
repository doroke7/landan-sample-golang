package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
)

const pubSubChannel = "channel:order.created"

func main() {
	// 收到中斷/終止訊號時 ctx 會被取消，下面的消費迴圈才會跟著結束，
	// 不是靠 process 被系統強制殺掉才停止。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1. 初始化標準 Redis 連線
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	// 2. 建立訂閱
	sub := rdb.Subscribe(ctx, pubSubChannel)
	defer sub.Close()

	// 💡 Subscribe 本身是非同步送出 SUBSCRIBE 指令，這裡先呼叫 Receive 確認訂閱
	// 真的成功建立，避免發布端在訂閱還沒完成前就搶先發送，導致第一批訊息收不到。
	if _, err := sub.Receive(ctx); err != nil {
		fmt.Printf("⚠️ 訂閱失敗: %v\n", err)
		return
	}

	msgCh := sub.Channel()

	fmt.Printf("🚪 [server] 已訂閱頻道 %s，等待訊息中（Ctrl+C 結束）...\n", pubSubChannel)

	// 3. 常駐消費訊息，直到 ctx 被取消（收到中斷訊號）或 channel 被關閉為止
	for {
		select {
		case <-ctx.Done():
			fmt.Println("👋 [server] 收到中斷訊號，結束消費")
			return
		case msg, bOk := <-msgCh:
			if !bOk {
				fmt.Println("👋 [server] channel 已關閉，退出消費迴圈")
				return
			}
			fmt.Printf("📩 [server] 頻道 %s 收到訊息：%s\n", msg.Channel, msg.Payload)
		}
	}
}
