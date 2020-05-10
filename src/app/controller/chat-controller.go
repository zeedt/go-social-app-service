package controller

import (
	"github.com/gin-gonic/gin"
	"go-social-app/src/app/models"
	"go-social-app/src/app/socket"
	"net/http"
)

func InitiateChatController(securedRoute *gin.RouterGroup)  {

	securedRoute.POST("/chat/send", func(context *gin.Context) {
		var chatMessage models.ChatMessage
		context.BindJSON(&chatMessage)
		username , _ := context.Get("oauth.credential")
		chat := models.Chat{
			Content:  chatMessage.Message,
			Sender:   username.(string),
			Receiver: chatMessage.To,
		}
		result := models.Db.Create(&chat)
		if result.Error != nil {
			context.JSON(http.StatusBadRequest, result.Error)
			return
		}

		context.String(http.StatusOK, "Done")
		socket.EmitToSocket(chat)

	})

}
