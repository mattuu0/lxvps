package lxdctl

import (
	"github.com/gorilla/websocket"
)

type LXShell struct {
	InstanceId string //インスタンスのID

	WsConn *websocket.Conn // フロントエンド用の websocket
	DataSocket *websocket.Conn // インスタンス側の stdin stdout用の Socket
	ControlSocket *websocket.Conn // インスタンス側の Control Websocket
}

// websocket のデータ構造
type WsData struct {
	Type string `json:"type"`
	Data string `json:"data"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}