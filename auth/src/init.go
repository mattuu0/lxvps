package main

import (
	"auth/controllers"
	"auth/middlewares"
	"html/template"
	"io"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (temp *TemplateRenderer) Render(writer io.Writer, name string, data interface{}, ctx echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = ctx.Echo().Reverse
	}

	return temp.templates.ExecuteTemplate(writer, name, data)
}

func SetupRouter(router *echo.Echo) {
	// logger 設定
	router.Use(middleware.Logger())

	// テンプレート
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	// レンダラー
	router.Renderer = renderer

	// ルーティング設定
	// ベーシックユーザーグループ
	// basicg := router.Group("/basic")
	// {
	// 	basicg.POST("/signup", controllers.CreateBasicUser)
	// 	basicg.POST("/login", controllers.LoginBasicUser)
	// }

	// React のビルド出力ディレクトリを指定
	buildDir := "dashboard"

	// サブパスの設定 (/_/ 配下に React アプリを配信)
	subPath := "/_"

	// 明示的なアセット (JS/CSS/画像など) の配信設定
	// 相対パスで配信するため、subPath 配下に設定
	router.Static(subPath+"/assets", filepath.Join(buildDir, "assets"))

	// subPath 配下のルートにアクセスがあった場合も index.html を返す
	router.GET(subPath, func(ctx echo.Context) error {
		indexPath := filepath.Join(buildDir, "index.html")
		return ctx.File(indexPath)
	})

	// SPA 対応: すべての subPath 配下の不明なルートを index.html にルーティング
	router.GET(subPath+"/*", func(ctx echo.Context) error {
		// index.html を返す
		indexPath := filepath.Join(buildDir, "index.html")
		return ctx.File(indexPath)
	})

	// アイコンフォルダを配信する
	router.Static("/assets", "./assets/icons")

	// info エンドポイント
	router.GET("/info/:userid", controllers.GetInfo)

	// 情報を取得する
	router.GET("/me", controllers.GetMe, middlewares.RequireAuth)

	// token を取得する
	router.GET("/token", controllers.GetToken, middlewares.RequireAuth)

	// アイコンを変更する
	router.POST("/icon", controllers.ChangeIcon, middlewares.RequireAuth)

	// アイコンを取得する
	router.GET("/icon/:userid",controllers.GetIcon)

	// ログアウト
	router.POST("/logout", controllers.Logout, middlewares.RequireAuth)

	// admin グループ
	adming := router.Group("/admin")
	{
		adming.POST("/signup", controllers.CreateAdminUser)
		adming.POST("/login", controllers.LoginAdminUser)
		adming.GET("/status", controllers.GetAdminStatus)
		adming.GET("/info", controllers.GetAdminInfo, middlewares.RequireAdminAuth)
		adming.POST("/logout", controllers.AdminLogout, middlewares.RequireAdminAuth)
	}

	// oauth グループ
	oauthg := router.Group("/oauth")
	{
		oauthg.GET("/:provider", controllers.StartOauth)
		oauthg.GET("/:provider/callback", controllers.CallbackOauth)
	}

	// api グループ
	apig := router.Group("/api")
	{
		// admin ミドルウェア設定
		apig.Use(middlewares.RequireAdminAuth)

		// ユーザーのグループ作成
		userg := apig.Group("/user")
		{
			// ユーザー一覧を取得する
			userg.GET("/all", controllers.GetAllUsers)

			// ユーザーを更新する
			userg.PUT("", controllers.UpdateUser)

			// ユーザーを削除する
			userg.DELETE("", controllers.DeleteOauth)

			// BAN を切り替える
			userg.PUT("/ban", controllers.ToggleBan)
		}

		// プロバイダグループ
		providerg := apig.Group("/providers")
		{
			// プロバイダー一覧を取得する
			providerg.GET("", controllers.GetProviders)

			// Oauth プロバイダ一覧取得
			providerg.GET("/oauth", controllers.GetOauthProviders)

			// Oauth プロバイダ一覧更新
			providerg.POST("/oauth", controllers.UpdateOauthProviders)

			// プロバイダー一覧を更新する
			providerg.POST("", controllers.UpdateProviders)

			// basic プロバイダ更新
			// providerg.PUT("/basic",controllers.BasicUpdate)
		}

		// ラベルグループを作る
		labelg := apig.Group("/labels")
		{
			// ラベルを取得する
			labelg.GET("", controllers.GetLabels)

			// ラベルを作成する
			labelg.POST("", controllers.CreateLabel)

			// ラベルを更新する
			labelg.PUT("", controllers.UpdateLabel)

			// ラベルを削除する
			labelg.DELETE("", controllers.DeleteLabel)
		}

		// セッショングループを作成する
		sessiong := apig.Group("/session")
		{
			// セッション一覧取得
			sessiong.GET("", controllers.GetSessions)

			// セッションを削除する
			sessiong.DELETE("", controllers.DeleteSession)
		}
	}
}
