package lxdctl

import (
	"github.com/gorilla/websocket"
)

type LXShell struct {
	InstanceId string //インスタンスのID

	WsConn *websocket.Conn // websocket 接続

	InstanceSocket *websocket.Conn // インスタンス側の websocket
}

// websocket のデータ構造
type WsData struct {
	Type string `json:"type"`
	Data string `json:"data"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}