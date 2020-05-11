package socket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"go-social-app/src/app/models"
	"log"
)

var ISocketServer *socketio.Server
var socketError error
var socketIds []string
var connectedUsersSocketMap = make(map[string]socketio.Conn)
var connectedUsersArray []models.SocketInfo

func InitiateSocketServer()  {
	ISocketServer, socketError = socketio.NewServer(nil)
	if socketError != nil {
		log.Fatal(socketError)
	}

	ISocketServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		ok := ISocketServer.JoinRoom("/", "my-room", s)
		fmt.Println("Joined ", ok)
		s.Emit("connection-received")
		return nil
	})

	ISocketServer.OnEvent("/", "join", func(s socketio.Conn, msg models.SocketInfo) {
		present := contains(socketIds, s.ID())
		if !present {
			connectedUsersSocketMap[s.ID()] = s
			msg.Name = msg.FirstName + " " + msg.LastName
			msg.ID = s.ID()
			socketIds = append(socketIds, s.ID())
			connectedUsersArray = append(connectedUsersArray, msg)
			ok := ISocketServer.BroadcastToRoom("/", "my-room","all-users", connectedUsersArray )
			fmt.Println("Broadcasting to room ", ok)

		}
	})
	ISocketServer.OnEvent("/", "bye", func(s socketio.Conn) string {
		fmt.Println("Bye bye")
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	ISocketServer.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	ISocketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		// Remove this socket info from each
		socketIds = removeSocket(s.ID(), socketIds)
		delete(connectedUsersSocketMap, s.ID())
		connectedUsersArray = updateConnectedUsersArray(s.ID())
		ISocketServer.BroadcastToRoom("/", "", "all-users", connectedUsersArray)

	})

	go ISocketServer.Serve()
	//defer ISocketServer.Close()
}


func contains(slice []string, val string) (bool) {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func removeSocket(value string, array []string ) []string {
	var newArray []string
	for index, item := range array{
		if item == value {
			newArray = append(newArray, array[(index+1) : len(array)]...)
			return newArray
		} else {
			newArray = append(newArray, item)
		}
	}
	return newArray
}

func updateConnectedUsersArray(socketId string) []models.SocketInfo {
	var newArray []models.SocketInfo
	for _, item := range connectedUsersArray{
		//if item.ID == socketId {
		//	newArray = append(newArray, connectedUsersArray[(index+1) : len(connectedUsersArray)]...)
		//	return newArray
		//} else {
		//	newArray = append(newArray, item)
		//}
		if item.ID != socketId {
			newArray = append(newArray, item)
		}
	}
	return newArray
}

func EmitToSocket(chat models.Chat)  {
	for _, value :=range connectedUsersArray {
		if value.Username == chat.Receiver || value.Username == chat.Sender {
			socket, found := connectedUsersSocketMap[value.ID]
			if found {
				socket.Emit("private-message-received", gin.H{
					"fromUsername": chat.Sender,
					"toUsername": chat.Receiver,
					"id":         chat.ID,
					"message" : chat.Content,
				})
			}
		}
	}
}