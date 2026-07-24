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
)

// client 端只用「一條」連線，靠 RequestId 讓多個並發呼叫共用它。
// 對照 sample/tcp_3（沒有 id、每次都是單一來回）才看得出多路復用多做了什麼，
// 關鍵動作都標了 // IMPORTANT 多路。

type result struct {
	Payload string
	Err     error
}

func main() {
	oConn, err := net.Dial("tcp", "127.0.0.1:9002")
	if err != nil {
		panic(err)
	}
	defer oConn.Close()

	oReader := bufio.NewReader(oConn)

	// IMPORTANT 多路：多個 goroutine 同時呼叫 call() 時，寫 request 也要序列化，
	// 理由跟 server 端的 oWriteMu 一樣——同一個 socket 不能被兩個 goroutine 同時寫。
	var oWriteMu sync.Mutex

	oPendingMu := sync.Mutex{}
	aPending := make(map[string]chan result)

	// IMPORTANT 多路：從連線建立到斷線為止，只有這一個 goroutine 負責讀這條連線，
	// 其他呼叫完全不碰 oReader。讀到 response 就靠 id 查表，把結果塞進對應呼叫方
	// 專屬的 channel——不共享讀取狀態，只讓一個人讀、其他人排隊拿結果。
	go func() {
		for {
			sId, sPayload, err := readMessage(oReader)
			if err != nil {
				return
			}

			oPendingMu.Lock()
			oChan, ok := aPending[sId]
			delete(aPending, sId)
			oPendingMu.Unlock()

			if ok {
				oChan <- result{Payload: sPayload}
			}
			// !ok：id 對不上任何還在等的呼叫方，直接丟掉。
		}
	}()

	call := func(sId string, sPayload string) (string, error) {
		oChan := make(chan result, 1)

		// IMPORTANT 多路：一定要「先登記、再送出」——在 Write 之前就把 channel
		// 放進 pending 表，確保 response 不管多快回來，讀取 goroutine 一定找得到
		// 對應的呼叫方。順序反過來的話，response 有可能在 channel 登記好之前
		// 就已經被讀取 goroutine 收到、查無此 id 而被直接丟掉。
		oPendingMu.Lock()
		aPending[sId] = oChan
		oPendingMu.Unlock()

		oWriteMu.Lock()
		err := writeMessage(oConn, sId, sPayload)
		oWriteMu.Unlock()
		if err != nil {
			return "", err
		}

		// IMPORTANT 多路：只 blocking 等自己專屬的 channel，不管其他並發中的
		// request 目前處理到哪、也不管 response 實際回來的順序，一定拿到
		// 屬於自己這一筆的結果。
		oResult := <-oChan
		return oResult.Payload, oResult.Err
	}

	fmt.Println("=== 多路復用測試：同時發出 3 筆 request（id=1,2,3） ===")
	fmt.Println("=== server 刻意讓 id=3 最快處理完、id=1 最慢，response 一定會亂序回來 ===")
	fmt.Println("=== 如果沒用 RequestId 配對，亂序回來的 response 會被誤認成別筆的結果 ===")

	var oWaitGroup sync.WaitGroup
	for i := 1; i <= 3; i++ {
		oWaitGroup.Add(1)
		go func(iId int) {
			defer oWaitGroup.Done()

			sId := strconv.Itoa(iId)
			sPayload, err := call(sId, "message-"+sId)
			if err != nil {
				fmt.Println("call failed:", err)
				return
			}

			fmt.Printf("client got response for id=%s: %s\n", sId, sPayload)
		}(i)
	}
	oWaitGroup.Wait()
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
