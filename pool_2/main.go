package main

import (
	"bufio"
	"errors"
	"net"
	"sync"
)

const (
	tcpMaxBodyLength = 1 << 12 // 4KB，避免錯誤/惡意的長度前綴把記憶體打爆
)

var (
	ErrTcpBodyTooLarge   = errors.New("tcp: body too large")
	ErrTcpMethodNotFound = errors.New("tcp: method not found")
	ErrTcpDialFailed     = errors.New("tcp: dial failed")
)

// tcpConn 是池子裡實際借還的單位：一條 net.Conn 綁一個專屬的 bufio.Reader。
// reader 要跟著連線的生命週期走、不能每次呼叫都重建（理由見 DecodeFrame 註解），
// 所以借出去、還回來都是這一整包一起借還，不會有連線和 reader 對不上的情況。
type TcpConn struct {
	conn   net.Conn
	reader *bufio.Reader
}

// 跟 TcpClient.NewTcpClient 的做法一致。
func NewTcpConn() *TcpConn {
	oConn, err := net.Dial("tcp", "127.0.0.1:8800")
	if err != nil {
		return nil
	}

	return &TcpConn{
		conn:   oConn,
		reader: bufio.NewReader(oConn),
	}
}

func (oSelf *TcpConn) Write(aBuf []byte) (int, error) {
	return oSelf.conn.Write(aBuf)
}

func (oSelf *TcpConn) Close() error {
	return oSelf.conn.Close()
}

// TcpPoolClient 用 sync.Pool 當連線池：不像 channel 版本有固定容量，
// sync.Pool 完全不設上限，Put 進去的東西也可能在任何一次 GC 時被悄悄清掉
// （不會呼叫 Close，連線就這樣被丟掉、底層 fd 靠 GC finalizer 或作業系統自己回收）。
// 換句話說 sync.Pool 天生是給「可以隨時重建、丟了也無所謂」的東西用的，
// 拿來裝有實體資源（TCP 連線）的物件，語意上比 channel 版本鬆散一些，這裡只是示範寫法。
type TcpPoolClient struct {
	pool sync.Pool
}

// NewTcpPoolClient 用 sync.Pool.New 接手「池子沒有現成連線時要怎麼生一條」，
// 不用像 channel 版本那樣在 get() 裡手動判斷空了要不要現撥。
func NewTcpPoolClient() *TcpPoolClient {
	return &TcpPoolClient{
		pool: sync.Pool{
			New: func() any {
				return NewTcpConn()
			},
		},
	}
}

// get 跟 sync.Pool 借一條連線；New 撥號失敗時會回傳 (*TcpConn)(nil)，
// 這裡要把它轉成 error，不然呼叫端會拿到一個看起來非 nil 介面、實際上是 nil 指標的 *TcpConn。
func (oSelf *TcpPoolClient) get() (*TcpConn, error) {
	oTcpConn, _ := oSelf.pool.Get().(*TcpConn)
	if oTcpConn == nil {
		return nil, ErrTcpDialFailed
	}
	return oTcpConn, nil
}

// put 把用完的連線還給 sync.Pool；沒有「池子滿了」這種狀態，
// 沒有 Close 的分支——sync.Pool 本身沒有提供「歸還時順便處理掉」的 hook。
func (oSelf *TcpPoolClient) put(oTcpConn *TcpConn) {
	oSelf.pool.Put(oTcpConn)
}
