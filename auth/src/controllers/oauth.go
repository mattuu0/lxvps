package controllers

import (
	"auth/logger"
	"auth/oauth2"
	"auth/services"
	"auth/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
)

func StartOauth(ctx echo.Context) error {
	provider := ctx.Param("provider")

	// oauth を更新
	oauth2.UseProviders()

	// IsMobile
	isMobile := ctx.QueryParam("ismobile")

	// popup 認証かどうか
	isPopup := ctx.QueryParam("popup")

	// 認証を開始
	oauth2.StartOauth(ctx, oauth2.OauthArgs{
		ProviderName: provider,
		IsMobile:     isMobile == "1",
		IsPopup:      isPopup == "1",
	})

	return nil
}

func CallbackOauth(ctx echo.Context) error {
	provider := ctx.Param("provider")

	// oauth を完了
	oauthResponse, err := oauth2.CallbackOauth(ctx, provider)

	// エラー処理
	if err != nil {
		// html を返す
		return utils.ErrorScreen(ctx, http.StatusInternalServerError, utils.GenID(), err, oauthResponse.IsPopup)
	}

	// ユーザー
	user := oauthResponse.User

	// popup
	isPopup := "0"
	if oauthResponse.IsPopup {
		isPopup = "1"
	}

	// ユーザーを作成
	token, err := services.LoginOauthUser(services.OauthUserArgs{
		Name:           GetName(user),
		Email:          user.Email,
		ProviderCode:   provider,
		ProviderUserID: user.UserID,
		RemoteIP:       ctx.RealIP(),
		UserAgent:      ctx.Request().UserAgent(),
		AvaterURL:      user.AvatarURL,
	})

	// エラー処理
	if err != nil {
		// return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		return utils.ErrorScreen(ctx, http.StatusInternalServerError, utils.GenID(), err, oauthResponse.IsPopup)
	}

	logger.Println(token)

	// キャッシュを無しにする
	ctx.Response().Header().Set("Expires", time.Unix(0, 0).Format(time.RFC1123))
	ctx.Response().Header().Set("Cache-Control", "no-cache, private, max-age=0")
	ctx.Response().Header().Set("Pragma", "no-cache")
	ctx.Response().Header().Set("X-Accel-Expires", "0")

	// モバイル場合
	if oauthResponse.IsMobile {
		return ctx.Redirect(http.StatusFound, "authbase://?token="+token)
	}

	return ctx.Render(http.StatusOK, "oauth-callback.html", echo.Map{"token": token, "isPopup": isPopup})
	// return ctx.JSON(http.StatusOK, echo.Map{"token": token})
	// return ctx.Redirect(http.StatusFound, "/auth/")
}

func GetName(user goth.User) string {
	result := ""

	logger.Println(user)

	if user.Name != "" {
		return user.Name
	}

	if user.NickName != "" {
		return user.NickName
	}

	if user.FirstName != "" {
		result = user.FirstName
	}

	if user.LastName != "" {
		if result != "" {
			result += " "
		}

		result += user.LastName
	}

	return result
}
