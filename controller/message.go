package controller

import (
	"fmt"
	"github.com/AdaQiao/simpleDouyin/model"
	"github.com/AdaQiao/simpleDouyin/service"
	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var tempChat = map[string][]model.Message{}

var MessageIdSequence = int64(1)

type ChatResponse struct {
	model.Response
	MessageList []model.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")

	if user, exist := UsersLoginInfo[token]; exist {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(user.Id, int64(userIdB))

		atomic.AddInt64(&MessageIdSequence, 1)
		curMessage := model.Message{
			Id:         MessageIdSequence,
			Content:    content,
			CreateTime: time.Now().Format(time.Kitchen),
		}

		if messages, exist := tempChat[chatKey]; exist {
			tempChat[chatKey] = append(messages, curMessage)
		} else {
			tempChat[chatKey] = []model.Message{curMessage}
		}
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")

	if user, exist := UsersLoginInfo[token]; exist {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(user.Id, int64(userIdB))
		//更新
		var err error
		var Chat []model.Message
		Chat, err = service.GetAllChat(user.Id, int64(userIdB))
		if err != nil {
			fmt.Sprintf("加载新聊天记录失败")
		} else {
			tempChat[chatKey] = Chat
		}
		c.JSON(http.StatusOK, ChatResponse{Response: model.Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "This user doesn't exist"})
	}
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
