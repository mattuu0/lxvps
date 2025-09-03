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
		Command:     []string{"su","-l"},
		WaitForWS:   true,
		Interactive: true,
		User: 0,
		Group: 0,
		Environment: map[string]string{
			"TERM" : "xterm-256color",
			"HOME" : "/root",
			"LANG" : "C.UTF-8",
			"USER" : "root",
		},
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

	logger.Println("インスタンス管理用のソケットを取得します")

	opAPI := op.Get()
	fds, ok := opAPI.Metadata["fds"].(map[string]interface{})
	if !ok {
		logger.Println("FDSメタデータを取得できません")
		return err
	}

	// データ用WebSocket（stdin/stdout）を取得
	var dataSecret string
	if secretVal, exists := fds["0"]; exists {
		dataSecret = secretVal.(string)
	}

	// 制御用WebSocketを取得
	var controlSecret string
	if secretVal, exists := fds["control"]; exists {
		controlSecret = secretVal.(string)
	}

	// エラー処理
	if dataSecret == "" || controlSecret == "" {
		logger.Println("必要なsecretを取得できません")
		return err
	}

	logger.Println("データ用WebSocketを取得します")
	
	// データ用WebSocket接続
	dataSocket, err := op.GetWebsocket(dataSecret)
	if err != nil {
		logger.Println("データWebSocket接続エラー:", err)
		return err
	}

	// データ用ソケットをセット
	ctrlShell.DataSocket = dataSocket

	logger.Println("制御用WebSocketを取得します")
	
	// 制御用WebSocket接続
	controlSocket, err := op.GetWebsocket(controlSecret)
	if err != nil {
		logger.Println("制御WebSocket接続エラー:", err)
		return err
	}

	// 制御用ソケットをセット
	ctrlShell.ControlSocket = controlSocket

	logger.PrintErr("ソケット取得完了")

	// プロセスが終わるまで待機
	err = op.Wait()
	if err != nil {
		logger.Println(err)
		return err
	}

	return nil
}
