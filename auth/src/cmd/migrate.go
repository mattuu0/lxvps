/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"auth/logger"
	"auth/models"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var Retry *bool

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "データベースをマイグレーションする",
	Long: `データベースをマイグレーションする`,
	Run: func(cmd *cobra.Command, args []string) {
		for {
			logger.Println("マイグレーションを実行しています...")

			// モデルを呼び出す
			err := models.Init()

			// エラー処理
			if err != nil {
				logger.PrintErr("マイグレーションに失敗しました", err)

				// 再試行するか
				if *Retry {
					logger.Println("500ms 後に再試行します")
					time.Sleep(500 * time.Millisecond)
					continue
				}

				os.Exit(1)
			}

			logger.Println("マイグレーション完了")
			break
		}

		// 終了する
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	Retry = migrateCmd.Flags().BoolP("retry", "r", false, "成功するまでリトライするか")
}
