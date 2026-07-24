package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 這份在 sample/tcp_3 的框架基礎上（4-byte 長度前綴解決黏包/拆包）加上 RequestId，
// 讓 client 可以在同一條連線上「同時」掛好幾筆還沒回應的 request——這就是多路復用
// （multiplexing）。跟 tcp_3 的差異都標了 // IMPORTANT 多路。
//
// message 格式很簡單，用一個 "|" 分隔："<id>|<payload>"，request/response 都一樣。

func main() {
	oListener, err := net.Listen("tcp", ":9002")
	if err != nil {
		panic(err)
	}
	defer oListener.Close()

	fmt.Println("server listening on :9002")

	for {
		oConn, err := oListener.Accept()
		if err != nil {
			continue
		}

		go handleConn(oConn)
	}
}

func handleConn(oConn net.Conn) {
	defer oConn.Close()

	oReader := bufio.NewReader(oConn)

	// IMPORTANT 多路：多個 goroutine 會同時想 Write 到同一個 net.Conn，寫入本身
	// 不是原子的，這裡只保護「寫一個完整 frame」這個動作；不然兩筆 response 的
	// bytes 可能交錯寫進去，變成誰都解不出來的垃圾。
	var oWriteMu sync.Mutex

	for {
		sId, sPayload, err := readMessage(oReader)
		if err != nil {
			return
		}

		// IMPORTANT 多路：收到 request 立刻丟進自己的 goroutine 處理，read 迴圈
		// 馬上回去讀下一筆。如果還是「處理完才讀下一筆」（跟 tcp_3 一樣），
		// client 就算並發送出好幾筆 request，server 這邊照樣是排隊處理，
		// 多路復用完全發揮不出來。
		go func(sId string, sPayload string) {
			// 故意讓每筆 request 的處理時間不一樣（id 數字越小睡越久），
			// 讓 response 確定會亂序回去——這樣才能真正驗證 client 端是
			// 靠 RequestId 配對回正確的呼叫方，不是「剛好照順序回來」的巧合。
			iId, _ := strconv.Atoi(sId)
			time.Sleep(time.Duration(400-iId*100) * time.Millisecond)

			fmt.Printf("server processed id=%s payload=%s\n", sId, sPayload)

			oWriteMu.Lock()
			defer oWriteMu.Unlock()

			// IMPORTANT 多路：response 一定要帶回「跟這個 request 一樣的 id」，
			// client 端全靠這個 id 才知道這筆 response 屬於哪個呼叫方；
			// 帶錯 id、或忘記帶，client 那邊的配對機制就直接失效。
			if err := writeMessage(oConn, sId, "echo:"+sPayload); err != nil {
				fmt.Println("server write failed:", err)
			}
		}(sId, sPayload)
	}
}

func encodeMessage(sId string, sPayload string) []byte {
	aBody := []byte(sId + "|" + sPayload)

	aFrame := make([]byte, 4+len(aBody))
	binary.BigEndian.PutUint32(aFrame[0:4], uint32(len(aBody)))
	copy(aFrame[4:], aBody)

	return aFrame
}

func writeMessage(oWriter io.Writer, sId string, sPayload string) error {
	_, err := oWriter.Write(encodeMessage(sId, sPayload))
	return err
}

// readMessage 先讀 4 byte 拿到長度、讀滿那個長度（解決黏包/拆包，跟 tcp_3 一樣），
// 再把內容用 "|" 切成 id 跟 payload 兩段。
func readMessage(oReader io.Reader) (sId string, sPayload string, err error) {
	aLengthBuf := make([]byte, 4)
	if _, err = io.ReadFull(oReader, aLengthBuf); err != nil {
		return "", "", err
	}

	iLength := binary.BigEndian.Uint32(aLengthBuf)

	aBody := make([]byte, iLength)
	if _, err = io.ReadFull(oReader, aBody); err != nil {
		return "", "", err
	}

	aParts := strings.SplitN(string(aBody), "|", 2)
	if len(aParts) != 2 {
		return "", string(aBody), nil
	}
	return aParts[0], aParts[1], nil
}
