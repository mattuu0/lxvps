package lxdctl

import (
	"app/logger"
	"io/ioutil"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
	"github.com/gorilla/websocket"
)

var (
	lxdClient lxd.InstanceServer // LXDクライアント
)

func Init() {
	// TLS証明書を読み込む
	clientCert, err := ioutil.ReadFile("./keys/client.crt")

	// エラー処理
	if err != nil {
		logger.Println(err)
		return
	}

	// TLS鍵を読み込む
	clientKey, err := ioutil.ReadFile("./keys/client.key")

	// エラー処理
	if err != nil {
		logger.Println(err)
		return
	}

	// LXDに接続
	server, err := lxd.ConnectLXD("https://192.168.100.20:8443", &lxd.ConnectionArgs{
		TLSClientCert:      string(clientCert),
		TLSClientKey:       string(clientKey),
		InsecureSkipVerify: true,
	})

	// エラー処理
	if err != nil {
		logger.Println("LXDに接続できません")
		logger.Println(err)
		return
	}

	// global にクライアントを設定
	lxdClient = server
}

// シェルを作成する
func InitLxd(InstanceId string, WsConn *websocket.Conn) error {
	// 管理用のクラスをインスタンス化する
	ctrlShell := LXShell{
		InstanceId: InstanceId,
		WsConn:     WsConn,
	}

	logger.Println("シェルを作成します")

	// Setup the exec request
	req := api.InstanceExecPost{
		Command:     []string{"bash"},
		WaitForWS:   true,
		Interactive: true,
		Width:       80,
		Height:      30,
	}

	// シェルを実行する
	logger.Println("シェルを実行します")

	// Setup the exec arguments (fds)
	args := lxd.InstanceExecArgs{
		Stdin:  &ctrlShell,
		Stdout: &ctrlShell,
		Stderr: &ctrlShell,
	}

	// Get the current state
	logger.Println("インスタンス内で実行します")

	// Get the current state
	op, err := lxdClient.ExecInstance(InstanceId, req, &args)
	if err != nil {
		logger.Println(err)
		return err
	}

	// オペレーションのメタデータを取得
	opAPI := op.Get()

	// WebSocketのsecretを取得
	var secret string
	if fds, ok := opAPI.Metadata["fds"].(map[string]interface{}); ok {
		if secretVal, exists := fds["control"]; exists {
			secret = secretVal.(string)
		}
	}

	// インスタンス側の websocket 取得
	instanceSocket, err := op.GetWebsocket(secret)

	// エラー処理
	if err != nil {
		logger.Println(err)
		return err
	}

	// インスタンス側の websocket を設定
	ctrlShell.InstanceSocket = instanceSocket

	// プロセスが終わるまで待機
	err = op.Wait()
	if err != nil {
		logger.Println(err)
		return err
	}

	return nil
}
