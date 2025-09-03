package lxdctl

import (
	"app/logger"
	"encoding/json"
	"strconv"

	"github.com/gorilla/websocket"
)

// writerも実装する (LXD -> フロントに転送用)
func (lxShell *LXShell) Write(data []byte) (n int, err error) {
	if lxShell.WsConn == nil {
		return 0, nil
	}

	// websocket にメッセージを送る
	err = lxShell.WsConn.WriteMessage(websocket.BinaryMessage, data)

	// エラー処理
	if err != nil {
		return 0, err
	}

	logger.Println("メッセージを送りました",data)

	return len(data), nil
}


// reader を実装する (フロント -> LXD に転送用)
func (lxShell *LXShell) Read(p []byte) (n int, err error) {
	// フロントから メッセージを受け取る
	_,readData,err := lxShell.WsConn.ReadMessage()

	// エラー処理
	if err != nil {
		logger.Println(err)
		return 0, err
	}

	// フロントからのメッセージ
	logger.Println("メッセージを受け取りました",readData)

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
		logger.Println("サイズを変更します",wsData.Cols,wsData.Rows)

		if lxShell.ControlSocket == nil {
			logger.Println("コントロール用のsocket がないです")
			return 0, nil
		}

		// サイズを変更する
		err := sendTermSize(lxShell.ControlSocket, wsData.Rows, wsData.Cols)

		// エラー処理
		if err != nil {
			logger.Println(err)
			return 0, err
		}

		logger.Println("サイズを変更しました")

		return 0, nil
	}


	// メッセージをコピーする
	copy(p, []byte(wsData.Data))
	return len(wsData.Data), nil
}

func sendTermSize(control *websocket.Conn, width, height int) error {
	// 公式WebUIと同じ形式のメッセージを作成
	msg := map[string]interface{}{
		"command": "window-resize",
		"args": map[string]string{
			"width":  strconv.Itoa(width),
			"height": strconv.Itoa(height),
		},
	}

	// JSONにシリアライズ
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	logger.Println("Sending resize message: %s", string(jsonData))

	// バイナリメッセージとして送信
	err = control.WriteMessage(websocket.BinaryMessage, jsonData)
	if err != nil {
		logger.Println("Failed to send resize message: %v", err)
		return err
	}

	logger.Println("Resize message sent successfully: %dx%d", width, height)
	return nil
}
