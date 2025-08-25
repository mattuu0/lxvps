package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/image/draw"
)

// SaveResizedImage は multipart.File からの画像を指定サイズにリサイズして保存します
// params:
//   - file: アップロードされたファイルのmultipart.File
//   - filename: ファイル名
//   - fileSize: ファイルサイズ
//   - width: リサイズ後の幅
//   - height: リサイズ後の高さ
//   - saveDir: 保存先ディレクトリ
//
// return:
//   - string: 保存されたファイルのパス
//   - error: エラー
func SaveResizedImage(file multipart.File, filename string, width, height int, saveDir string) (string, error) {

	// ファイル名の拡張子をチェック
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		// 拡張子がない場合はPNGとして扱う
		filename = filename + ".png"
	} else {
		// 出力は常にPNGなので、元の拡張子をPNGに置き換える
		filename = strings.TrimSuffix(filename, ext) + ".png"
	}

	// ファイルポインタを先頭に戻す
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to seek file: %w", err)
	}

	// 画像をデコード (EXIFデータは除去される)
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// リサイズした画像を保存
	return saveResizedImageToPNG(img, filename, width, height, saveDir)
}

// SaveResizedImageFromBytes はバイトスライスからの画像を指定サイズにリサイズして保存します
// params:
//   - imageData: 画像データのバイトスライス
//   - filename: ファイル名
//   - width: リサイズ後の幅
//   - height: リサイズ後の高さ
//   - saveDir: 保存先ディレクトリ
//
// return:
//   - string: 保存されたファイルのパス
//   - error: エラー
func SaveResizedImageFromBytes(imageData []byte, filename string, width, height int, saveDir string) (string, error) {
	if len(imageData) == 0 {
		return "", errors.New("empty image data")
	}

	// ファイル名の拡張子をチェック
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		// 拡張子がない場合はPNGとして扱う
		filename = filename + ".png"
	} else {
		// 出力は常にPNGなので、元の拡張子をPNGに置き換える
		filename = strings.TrimSuffix(filename, ext) + ".png"
	}

	// バイトスライスをReaderに変換
	reader := bytes.NewReader(imageData)

	// 画像をデコード
	img, _, err := image.Decode(reader)
	if err != nil {
		return "", fmt.Errorf("failed to decode image from bytes: %w", err)
	}

	// リサイズした画像を保存
	return saveResizedImageToPNG(img, filename, width, height, saveDir)
}

// SaveResizedImageFromURL はURLから画像をダウンロードしてリサイズし保存します
// params:
//   - imageURL: 画像のURL
//   - filename: 保存するファイル名（指定がない場合はURLから取得）
//   - width: リサイズ後の幅
//   - height: リサイズ後の高さ
//   - saveDir: 保存先ディレクトリ
//
// return:
//   - string: 保存されたファイルのパス
//   - error: エラー
func SaveResizedImageFromURL(imageURL string, filename string, width, height int, saveDir string) (string, error) {
	// URLが空でないか確認
	if imageURL == "" {
		return "", errors.New("empty image URL")
	}

	// ファイル名が指定されていない場合はURLから取得
	if filename == "" {
		urlParts := strings.Split(imageURL, "/")
		if len(urlParts) > 0 {
			filename = urlParts[len(urlParts)-1]
			// クエリパラメータを除去
			if idx := strings.Index(filename, "?"); idx != -1 {
				filename = filename[:idx]
			}
		}

		// それでもファイル名がない場合はタイムスタンプを使用
		if filename == "" {
			filename = fmt.Sprintf("image_%d.png", time.Now().UnixNano())
		}
	}

	// ファイル名の拡張子をチェック
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		// 拡張子がない場合はPNGとして扱う
		filename = filename + ".png"
	} else {
		// 出力は常にPNGなので、元の拡張子をPNGに置き換える
		filename = strings.TrimSuffix(filename, ext) + ".png"
	}

	// HTTPリクエストを作成
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image from URL: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスのステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	// 画像をデコード
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to decode image from URL: %w", err)
	}

	// リサイズした画像を保存
	return saveResizedImageToPNG(img, filename, width, height, saveDir)
}

// saveResizedImageToPNG は画像を指定サイズにリサイズしてPNG形式で保存する共通処理
// params:
//   - img: リサイズ対象の画像
//   - filename: 保存するファイル名
//   - width: リサイズ後の幅
//   - height: リサイズ後の高さ
//   - saveDir: 保存先ディレクトリ
//
// return:
//   - string: 保存されたファイルのパス
//   - error: エラー
func saveResizedImageToPNG(img image.Image, filename string, width, height int, saveDir string) (string, error) {
	// リサイズ先の画像を作成（アスペクト比を保持せずに強制的にリサイズ）
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	// リサイズ実行
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	// 保存先ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 保存先パスを作成
	savePath := filepath.Join(saveDir, filepath.Base(filename))

	// ファイルを作成
	outFile, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// PNG形式で画像を保存
	if err := png.Encode(outFile, dst); err != nil {
		return "", fmt.Errorf("failed to encode image as PNG: %w", err)
	}

	return savePath, nil
}

