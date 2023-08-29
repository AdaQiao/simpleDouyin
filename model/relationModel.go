package model

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type FollowActionMessage struct {
	Token      string
	ToUserId   int64
	ActionType int32
}
