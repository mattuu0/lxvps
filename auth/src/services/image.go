package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	_ "image/jpeg"
	_ "image/gif"
	
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/image/draw"
)

// 定数
const (
	PngPrefix  = "data:image/png;base64,"
	JpegPrefix = "data:image/jpeg;base64,"
	JpgPrefix  = "data:image/jpg;base64,"
	GifPrefix  = "data:image/gif;base64,"

	// アイコンのディレクトリ
	IconDir = "./assets/icons"

	IconWidth  = 256 // 画像サイズ
	IconHeight = 256 // 画像サイズ

	MaxImageSize = 10485760 // 10MB
)

// ProcessAndSaveImage は文字列データをデコードし、画像かを検証し、256x256にリサイズしてPNGで保存します
func ProcessAndSaveImage(filePath string, dataString string, maxDimension int) error {
	// データの形式を判定
	var encodedData string
	var format string

	switch {
	case strings.HasPrefix(dataString, PngPrefix):
		encodedData = strings.TrimPrefix(dataString, PngPrefix)
		format = "png"
	case strings.HasPrefix(dataString, JpegPrefix):
		encodedData = strings.TrimPrefix(dataString, JpegPrefix)
		format = "jpeg"
	case strings.HasPrefix(dataString, JpgPrefix):
		encodedData = strings.TrimPrefix(dataString, JpgPrefix)
		format = "jpeg"
	case strings.HasPrefix(dataString, GifPrefix):
		encodedData = strings.TrimPrefix(dataString, GifPrefix)
		format = "gif"
	default:
		return errors.New("無効なデータ形式: サポートされている画像のBase64エンコードデータではありません")
	}

	// Base64デコード
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return fmt.Errorf("Base64デコードエラー: %v", err)
	}

	// 画像として解析
	img, imgFormat, err := image.Decode(bytes.NewReader(decodedData))
	if err != nil {
		return fmt.Errorf("画像デコードエラー: %v", err)
	}

	// 画像フォーマットの確認（追加の検証）
	if imgFormat != format && !(imgFormat == "jpeg" && (format == "jpg" || format == "jpeg")) {
		return fmt.Errorf("データ形式が不一致: プレフィックスは %s だが実際の画像は %s", format, imgFormat)
	}

	// 画像サイズの検証
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width > maxDimension || height > maxDimension {
		return fmt.Errorf("画像サイズが大きすぎます: %dx%d (最大許容サイズ: %dx%d)", width, height, maxDimension, maxDimension)
	}

	// 256x256にリサイズ
	dst := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.CatmullRom.Scale(dst, dst.Rect, img, bounds, draw.Over, nil)

	// 出力ファイルを作成
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ファイル作成エラー: %v", err)
	}
	defer outFile.Close()

	// PNGとして保存
	err = png.Encode(outFile, dst)
	if err != nil {
		return fmt.Errorf("PNG エンコードエラー: %v", err)
	}

	return nil
}

// ProcessImageFromURL はURLから画像を取得し、検証し、256x256にリサイズしてPNGで保存します
func ProcessImageFromURL(filePath string, imageURL string, maxDimension int, timeoutSeconds int) error {
	// HTTPクライアントの設定（タイムアウト付き）
	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}

	// URLからのリクエスト
	resp, err := client.Get(imageURL)
	if err != nil {
		return fmt.Errorf("URLからの画像取得エラー: %v", err)
	}
	defer resp.Body.Close()

	// ステータスコードの確認
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP エラー: %s", resp.Status)
	}

	// Content-Typeの確認
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf("無効なコンテンツタイプ: %s", contentType)
	}

	// レスポンスボディを読み込む
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("画像データの読み込みエラー: %v", err)
	}

	// 画像として解析
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return fmt.Errorf("画像デコードエラー: %v", err)
	}

	// 画像サイズの検証
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width > maxDimension || height > maxDimension {
		return fmt.Errorf("画像サイズが大きすぎます: %dx%d (最大許容サイズ: %dx%d)", width, height, maxDimension, maxDimension)
	}

	// 256x256にリサイズ
	dst := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.CatmullRom.Scale(dst, dst.Rect, img, bounds, draw.Over, nil)

	// 出力ファイルを作成
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ファイル作成エラー: %v", err)
	}
	defer outFile.Close()

	// PNGとして保存
	err = png.Encode(outFile, dst)
	if err != nil {
		return fmt.Errorf("PNG エンコードエラー: %v", err)
	}

	return nil
}
