package lxdctl

import (
	"app/logger"
	"encoding/json"
	"strconv"

	"github.com/canonical/lxd/shared/api"
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
	_,readData,err := lxShell.WsConn.ReadMessage()

	// エラー処理
	if err != nil {
		logger.Println(err)
		return 0, err
	}

	// json構造体
	wsData := WsData{}

	// jsonデコードする
	err = json.Unmarshal(readData, &wsData)

	// エラー処理
	if err != nil {
		logger.Println(err)
		return 0, err
	}

	logger.Println("メッセージを受け取りました",wsData)

	// もし データが resize の時
	if wsData.Type == "resize" {
		if lxShell.InstanceSocket == nil {
			return 0, nil
		}

		// サイズを変更する
		err := sendTermSize(lxShell.InstanceSocket, wsData.Cols, wsData.Rows)

		// エラー処理
		if err != nil {
			logger.Println(err)
			return 0, err
		}

		return 0, nil
	}


	// メッセージをコピーする
	copy(p, []byte(wsData.Data))

	return len(wsData.Data), nil
}

func sendTermSize(control *websocket.Conn, width, height int) error {
	msg := api.InstanceExecControl{}
	msg.Command = "window-resize"
	msg.Args = make(map[string]string)
	msg.Args["width"] = strconv.Itoa(width)
	msg.Args["height"] = strconv.Itoa(height)

	return control.WriteJSON(msg)
}