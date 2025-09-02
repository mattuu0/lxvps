package services

import (
	"app/logger"
	"io/ioutil"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
)

func TestLxd() {
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

	// インスタンスを作成する
	req := api.InstancesPost{
		Name: "test-instance",
		Source: api.InstanceSource{
			Type:     "image",
			Alias:    "22.04",
			Mode:     "pull",
			Protocol: "simplestreams",
			Server:   "https://cloud-images.ubuntu.com/minimal/releases/",
		},
		Type: "container",
	}

	// インスタンスを作成
	operate, err := server.CreateInstance(req)

	// エラー処理
	if err != nil {
		logger.Println("インスタンスを作成できません")
		logger.Println(err)
		return
	}

	// インスタンス作成を待機する
	logger.Println("インスタンスを作成しています")
	err = operate.Wait()

	// エラー処理
	if err != nil {
		logger.Println("インスタンスを作成できません")
		logger.Println(err)
		// return
	}

	// インスタンスを開始する
	// Get LXD to start the instance (background operation)
	reqState := api.InstanceStatePut{
		Action:  "start",
		Timeout: -1,
	}

	// Start the instance
	op, err := server.UpdateInstanceState("test-instance", reqState, "")
	if err != nil {
		logger.Println("インスタンスを開始できません")
		logger.Println(err)
		return
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		logger.Println("インスタンスを開始できません")
		logger.Println(err)
		// return
	}
}
