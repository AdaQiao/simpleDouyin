package controller

import (
	"github.com/AdaQiao/simpleDouyin/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := UsersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: []model.User{model.DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: []model.User{model.DemoUser},
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: []model.User{model.DemoUser},
	})
}
