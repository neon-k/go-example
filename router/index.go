package router

import (
	handler "github.com/api/handler/rest"
	"github.com/api/infra/persistence"
	"github.com/api/middleware"
	"github.com/api/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/api/docs"
	"github.com/gin-gonic/gin"
)

// @title APIドキュメントのタイトル
// @version バージョン(1.0)
// @description 仕様書に関する内容説明
// @termsOfService 仕様書使用する際の注意事項

// @contact.name APIサポーター
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name ライセンス(必須)
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:4000
// @BasePath /
func StartRouter() *gin.Engine {
	engine := gin.Default()

	engine.Use(middleware.SetDB)

	// =====================================
	// 依存関係を注入
	// =====================================

	// example
	todoPersistence := persistence.NewTodoPersistence()
	todoUseCase := usecase.NewTodoUseCase(todoPersistence)
	todoHandler := handler.NewTodokHandler(todoUseCase)

	// user
	userPersistence := persistence.NewUserPersistence()
	userUseCase := usecase.NewUserseCase(userPersistence)
	userHandler := handler.NewUserkHandler(userUseCase)

	// post
	postPersistence := persistence.NewPostPersistence()
	potsUseCase := usecase.NewPostCase(postPersistence)
	postHandler := handler.NewPostHandler(potsUseCase)

	// good
	goodPersistence := persistence.NewGoodPersistence()
	goodUseCase := usecase.NewGoodCase(goodPersistence)
	goodHandler := handler.NewGoodHandler(goodUseCase)

	// refreshToken
	refreshTokenPersistence := persistence.NewRefreshTokenPersistence()
	refreshTokenUseCase := usecase.NewRefreshTokenCase(refreshTokenPersistence)

	// jwt
	jwtHandler := handler.NewJwtHandler(userUseCase, refreshTokenUseCase)

	// =====================================
	// ルーティング
	// =====================================

	api := engine.Group("api")
	{
		v1 := api.Group("v1")
		{
			v1.GET("/ping", handler.ApiPing)

			// user
			v1.GET("/users", userHandler.GetUserAll)
			v1.POST("/signup", userHandler.AddUser)

			// example
			v1.GET("/example", todoHandler.Index)

			// example
			v1.GET("/firestore", handler.firestore)

			// jwt
			v1.POST("/login", jwtHandler.AuthMiddleware().LoginHandler)

			v1.Use(jwtHandler.AuthMiddleware().MiddlewareFunc())
			{

				// user
				v1.DELETE("/user/:id", userHandler.DeleteUser)

				// refreshToken
				v1.PATCH("/refresh_token", jwtHandler.RefreshToken)

				// posts
				v1.GET("/posts", postHandler.GetPostAll)
				v1.GET("/posts/:id", postHandler.GetPostAll)

				// good
				v1.POST("/good/:id", goodHandler.SetGood)

				user := v1.Group("user")
				{
					user.GET("/posts/:id", postHandler.GetUserPosts)
				}

				self := v1.Group("self")
				{
					// user
					self.GET("/user", userHandler.GetCurrentUser)
					self.PATCH("/user", userHandler.UpdateUser)

					// post
					self.POST("/post", postHandler.AddPost)
					self.GET("/posts", postHandler.GetCurrentPosts)
					self.PATCH("/post/:id", postHandler.UpdatePost)
					self.DELETE("/post/:id", postHandler.DeletePost)
				}
			}
		}
	}

	url := ginSwagger.URL("/swagger/doc.json")
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return engine
}
