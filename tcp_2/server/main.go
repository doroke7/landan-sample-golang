package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// 這份是 sample/tcp_2 的複製，多加一個 countingReader 來「證明」黏包/拆包真的發生了，
// 不是紙上談兵：每讀完一筆完整訊息，就印出這期間底層 conn.Read 實際被呼叫了幾次。
//
//   - 黏包：client 一次把 3 筆訊息的 bytes 串在一起送出去（見 client/main.go），
//     server 讀第 1 筆訊息時，底層 Read 很可能一次就把 3 筆的 bytes 全部收進
//     bufio.Reader 的內部緩衝區——所以第 2、3 筆訊息會顯示「底層 Read 呼叫 0 次」，
//     因為資料早就在緩衝區裡了，根本不需要再問 kernel 要新資料。
//   - 拆包：client 送一筆超大訊息，遠大於 bufio.Reader 一次能裝的量，
//     底層 Read 勢必要被呼叫很多次才能把這一筆訊息收滿。

func main() {
	oListener, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err)
	}
	defer oListener.Close()

	fmt.Println("server listening on :9001")

	for {
		oConn, err := oListener.Accept()
		if err != nil {
			continue
		}

		go handleConn(oConn)
	}
}

// countingReader 包一層在 net.Conn 外面，每次底層 Read 被呼叫就記一次數，
// 用來觀察「一筆邏輯訊息」實際上對應了幾次「實體的 TCP Read」。
type CountingReader struct {
	reader io.Reader // CountingReader 繼承 了  io.Reader
	Count  int
}

func (oSelf *CountingReader) Read(aBuf []byte) (int, error) {
	oSelf.Count++
	return oSelf.reader.Read(aBuf)
}

func handleConn(oConn net.Conn) {
	defer oConn.Close()

	oCounterReader := &CountingReader{reader: oConn}
	oReader := bufio.NewReader(oCounterReader)

	iIndex := 0
	for {
		iCountBefore := oCounterReader.Count

		/*
			    IMPORTANT
				原本的機制是， 一個同步的讀取resquest，然後再同步的寫 response。
				換言之這個response 一定確切被這個 resquest 相關的數據後寫入。


		*/

		sMessage, err := readMessage(oReader)
		if err != nil {
			return
		}
		iIndex++

		sPreview := sMessage
		if len(sPreview) > 20 {
			sPreview = sPreview[:20] + "..."
		}

		fmt.Printf(
			"server received #%d: len=%d preview=%q (這筆訊息期間底層 Read 被呼叫了 %d 次)\n",
			iIndex, len(sMessage), sPreview, oCounterReader.Count-iCountBefore,
		)

		sResponse := fmt.Sprintf("echo #%d len=%d", iIndex, len(sMessage))
		if err := writeMessage(oConn, sResponse); err != nil {
			return
		}
	}
}

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

// readMessage 先讀 4 byte 拿到長度，再讀滿那個長度——不管底層一次送幾個 byte 過來，
// 讀滿了才算一個完整的 message，這就是解決黏包/拆包的關鍵。
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
