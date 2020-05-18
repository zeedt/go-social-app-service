package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/maxzerbini/oauth"
	"go-social-app/src/app/auth"
	"go-social-app/src/app/controller"
	"go-social-app/src/app/models"
	_ "go-social-app/src/app/models"
	"go-social-app/src/app/socket"
	"net/http"
	"time"
)

func main() {
	models.Migrate()
	route := gin.Default()
	route.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Access-Control-Request-Headers", "Accept", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		//AllowAllOrigins:  true,
		AllowOrigins: []string{ "http://localhost:3000","http://localhost:5000"},
		//AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           86400,
	}))
	route.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Yes")
	})

	socket.InitiateSocketServer()
	route.GET("/socket.io/*any", gin.WrapH(socket.ISocketServer))
	route.POST("/socket.io/*any", gin.WrapH(socket.ISocketServer))

	securedRoute := getSecurededRoute(route)

	controller.InitiateUserController(route, securedRoute)
	controller.InitiatePostRoute(securedRoute)
	controller.InitiateChatController(securedRoute)
	controller.InitiateCommentRoute(securedRoute)

	route.Run(":3004")
}

func getSecurededRoute(route *gin.Engine) *gin.RouterGroup {
	oauthBearerServer := oauth.NewOAuthBearerServer("random-string", 120 * time.Minute, &auth.OauthVerifier{}, nil)

	route.POST("oauth/token", oauthBearerServer.UserCredentials)
	route.POST("client/token", oauthBearerServer.ClientCredentials)
	authorized := route.Group("/")
	authorized.Use(oauth.Authorize("random-string", nil))
	return authorized

}