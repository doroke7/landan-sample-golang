package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	oConn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		panic(err)
	}
	defer oConn.Close()

	oReader := bufio.NewReader(oConn)

	fmt.Println("=== 黏包測試：3 筆小 message 串成一個 []byte，一次 Write 出去 ===")
	pack(oConn, oReader)

	fmt.Println()
	fmt.Println("=== 拆包測試：送 1 筆超大 message，底層一定要分好幾次 Read 才收得完 ===")
	unpack(oConn, oReader)
}

func pack(oConn net.Conn, oReader *bufio.Reader) {
	aMessages := []string{"msg-A", "msg-B", "msg-C"}

	// 故意把 3 筆訊息各自組好 frame 之後「串在同一個 []byte 裡，一次 Write 出去」，
	// 而不是分開呼叫 3 次 Write。這樣可以確保底層一定是當成一坨連續的 bytes 送出去，
	// server 端第一次 Read 很可能就把 3 筆訊息的 bytes 全部收進來——這就是黏包。
	var aPayload []byte
	for _, sMessage := range aMessages {
		aPayload = append(aPayload, encodeMessage(sMessage)...)
	}

	if _, err := oConn.Write(aPayload); err != nil {
		panic(err)
	}

	// 3 筆 request 對應 3 筆 response，依序讀回來，驗證每一筆都完整、沒有互相混到。
	for range aMessages {
		sResponse, err := readMessage(oReader)
		if err != nil {
			panic(err)
		}
		fmt.Println("client received:", sResponse)
	}
}

func unpack(oConn net.Conn, oReader *bufio.Reader) {
	// 5MB 的內容，遠大於 bufio.Reader 預設緩衝區（4096 bytes）能一次裝的量，
	// 底層 Read 勢必要被呼叫很多次才能把這一筆訊息收滿——這就是拆包。
	sBigMessage := strings.Repeat("A", 5*1024*1024)

	if err := writeMessage(oConn, sBigMessage); err != nil {
		panic(err)
	}

	sResponse, err := readMessage(oReader)
	if err != nil {
		panic(err)
	}
	fmt.Println("client received:", sResponse)
}

// encodeMessage/writeMessage/readMessage 跟 server 那份一模一樣——這是故意的，
// 兩邊本來就要照同一套規則組包/拆包才讀得懂彼此，不是誰依賴誰。

func encodeMessage(sMessage string) []byte {
	aBody := []byte(sMessage)

	aFrame := make([]byte, 4+len(aBody))
	binary.BigEndian.PutUint32(aFrame[0:4], uint32(len(aBody)))
	copy(aFrame[4:], aBody)

	return aFrame
}

func writeMessage(oWriter io.Writer, sMessage string) error {
	_, err := oWriter.Write(encodeMessage(sMessage))
	return err
}

func readMessage(oReader io.Reader) (string, error) {
	aLengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(oReader, aLengthBuf); err != nil {
		return "", err
	}

	iLength := binary.BigEndian.Uint32(aLengthBuf)

	aBody := make([]byte, iLength)
	if _, err := io.ReadFull(oReader, aBody); err != nil {
		return "", err
	}

	return string(aBody), nil
}
