package model

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

type CommentActionRequest struct {
	Token       string
	VideoId     int64
	ActionType  int32
	CommentText string
	CommentId   int64
}
