package lxdctl

import (
	"app/logger"

	"github.com/gorilla/websocket"
)

// writerも実装する
func (lxShell *LXShell) Write(data []byte) (n int, err error) {
	// websocket にメッセージを送る
	err = lxShell.WsConn.WriteMessage(websocket.BinaryMessage, data)

	// エラー処理
	if err != nil {
		return 0, err
	}

	logger.Println("メッセージを送りました",data)

	return len(data), nil
}


// reader を実装する
func (lxShell *LXShell) Read(p []byte) (n int, err error) {
	// websocket からメッセージを受け取る
	_,data,err := lxShell.WsConn.ReadMessage()

	// エラー処理
	if err != nil {
		logger.Println(err)
		return 0, err
	}

	logger.Println("メッセージを受け取りました",data)

	// メッセージをコピーする
	copy(p, data)

	return len(data), nil
}
