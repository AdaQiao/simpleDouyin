package model

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
type UploadViewReq struct {
	token    string
	viewUrl  string
	coverUrl string
}
