package model

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
type uploadViewReq struct {
	token    string
	viewUrl  string
	coverUrl string
}
