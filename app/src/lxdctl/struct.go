package lxdctl

import (
	"github.com/gorilla/websocket"
)

type LXShell struct {
	InstanceId string //インスタンスのID

	WsConn *websocket.Conn // websocket 接続
}
