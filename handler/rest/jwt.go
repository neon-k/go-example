package rest

import (
	"net/http"
	"os"
	"time"

	"github.com/api/domain/model"
	"github.com/api/usecase"
	"github.com/api/util"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/theckman/go-securerandom"
)

type TLoginReq struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type TLoginRes struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	Expire       time.Time `json:"expire"`
}

type TRefreshToken struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}

type TRefreshTokenReq struct {
	Token string `json:"token" validate:"required"`
}

type JwtHandler interface {
	AuthMiddleware() *jwt.GinJWTMiddleware
	RefreshToken(c *gin.Context)
}

// usercaseのintefaceと紐ずける
type jwtHandler struct {
	userUseCase         usecase.UserUseCase
	refreshTokenUseCase usecase.RefreshTokenUseCase
}

// NewTodoUseCase : Todo データに関する Handler を生成
func NewJwtHandler(uu usecase.UserUseCase, ur usecase.RefreshTokenUseCase) JwtHandler {
	return &jwtHandler{
		userUseCase:         uu,
		refreshTokenUseCase: ur,
	}
}

// リフレッシュトークン
var maxRefresh = time.Hour * 24 * 30 * 6

// @description ユーザー認証
// @version 1.0
// @Tags auth
// @Summary ユーザー認証
// @param request body TLoginReq false "リクエスト"
// @accept application/x-json-stream
// @Success 200 {object} TLoginRes
// @router /api/v1/login [POST]
func (jh jwtHandler) AuthMiddleware() *jwt.GinJWTMiddleware {
	var authMiddleware, _ = jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "chats",
		Key:        []byte(os.Getenv("JWT_SECRET_KEY")),
		Timeout:    time.Hour * 24 * 30,
		MaxRefresh: maxRefresh,

		// ===========================================
		// ログイン時
		// ===========================================

		// ログイン時呼ばれる関数
		// 一番はじめにここに入る
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// dbを取得
			DB := util.DB(c)

			var loginVal TLoginReq

			// reqのデータを取得
			if err, _ := util.GetRequestValidate(c, &loginVal); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			// apiからきたemail、パスワードを取得
			email := loginVal.Email
			password := loginVal.Password

			// ユーザー情報を取得
			user, errDB := jh.userUseCase.GetCurrentUser(DB, password, email)

			if errDB != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			// 有効期限を発行
			now := time.Now()
			expire := now.Add(maxRefresh)

			// ランダムのバイト数を生成（リフレッシュトークンになる）
			rStr, errRStr := securerandom.Base64OfBytes(64)

			if errRStr != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			// リフレッシュトークンをDBに格納
			if errRefreshToken := jh.refreshTokenUseCase.AddRefreshToken(DB, rStr, &expire); errRefreshToken != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			// refreshTokenデータに格納
			c.Set("refreshToken", &TRefreshToken{
				Token:  rStr,
				Expire: expire,
			})

			// 取得したユーザーを返す
			return user, nil
		},

		// ログイン時に呼ばれる関数
		// tokenにデータを詰め込む
		// Authenticator
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {

				claims := jwt.MapClaims{
					"userID": v.ID,
					"name":   v.Name,
				}

				return claims
			}
			return jwt.MapClaims{}
		},

		// ログイン時に返すres
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			// dbを取得
			refreshToken := c.MustGet("refreshToken").(*TRefreshToken)

			c.JSON(code, gin.H{
				"token":        token,
				"refreshToken": &refreshToken.Token,
				"expire":       expire.Format(time.RFC3339),
			})
		},

		// ===========================================
		// Middleware (token確認 +  tokenのデータを取得)
		// ===========================================

		// MiddlewareFuncを使うと呼ばれる
		// tokenの中身を確認、idを取得してユーザー情報があるか確認する
		IdentityHandler: func(c *gin.Context) interface{} {

			// dbを取得
			DB := util.DB(c)

			// tokenの中身を確認
			claims := jwt.ExtractClaims(c)

			// ユーザーIDを取得
			id := claims["userID"].(float64)

			// ユーザー情報を取得
			user, errDB := jh.userUseCase.GetCurrentUserID(DB, id)

			if errDB != nil {
				return nil
			}

			// ログインしたユーザー情報を返す
			return user
		},

		// MiddlewareFuncを使うと呼ばれる。
		// IdentityHandlerで返された値が入る
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// ログインしたユーザー情報があるか確認
			if v, ok := data.(*model.User); ok {

				// ログインしたユーザー情報をginに格納
				util.SetCurrentUser(c, v)
				return true
			}

			return false
		},

		// ===========================================
		// エラー時に呼ばれる
		// ===========================================

		// エラーした時に呼ばれる
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		TokenLookup:   "header: Authorization, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	return authMiddleware
}

// @description tokenを再発行
// @version 1.0
// @Tags refresh_token
// @Summary tokenを再発行
// @param request body TRefreshTokenReq false "リクエスト"
// @accept application/x-json-stream
// @Security ApiKeyAuth
// @Success 200 {object} TRefreshToken
// @router /api/v1/refresh_token [PATCH]
func (jh jwtHandler) RefreshToken(c *gin.Context) {
	// 型を定義
	var req TRefreshTokenReq

	// DBを定義
	DB := util.DB(c)

	// reqを取得
	if err, errorMessages := util.GetRequestValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err, "messages": errorMessages})
		return
	}

	// リフレッシュトークンをチェック
	if err := jh.refreshTokenUseCase.CheckRefreshToken(DB, req.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"RefreshToken": gin.H{"text": "refreshToken does not match", "tag": "notmatch"}})
		return
	}

	// トークンを再発行
	jh.AuthMiddleware().RefreshHandler(c)
}
