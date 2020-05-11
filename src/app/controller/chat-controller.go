package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-social-app/src/app/models"
	"go-social-app/src/app/socket"
	"net/http"
	"sort"
	"strconv"
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

	securedRoute.GET("/chat/load-with-less-id/:id/:otherUser", func(context *gin.Context) {
		var id int
		passedId, found := context.Params.Get("id")
		otherUser, otherUserfound := context.Params.Get("otherUser")
		if !found {
			context.JSON(http.StatusBadRequest, gin.H{"message":"ID not passed"})
			return
		} else {
			id, _ = strconv.Atoi(passedId)
		}

		if !otherUserfound {
			context.JSON(http.StatusBadRequest, gin.H{"message":"Username not passed"})
			return
		}
		fmt.Println(id)
		username , _ := context.Get("oauth.credential")
		var chats []models.Chat

		if id < 1 {
			models.Db.Order("id desc").Where("receiver = ?  and sender = ? ", otherUser, username.(string)).Or("receiver = ?  and sender = ? ",
				username.(string), otherUser).Limit(10).Find(&chats)
			sort.Sort(models.ChatsById(chats))
		} else {
			models.Db.Order("id asc").Where("receiver = ?  and sender = ? and id < ? ", otherUser, username.(string), id).Or("receiver = ?  and sender = ? and id < ?",
				username.(string), otherUser, id).Limit(10).Find(&chats)
		}

		context.JSON(http.StatusOK, chats)
	})

}
