package model

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
type UploadViewReq struct {
	Token    string
	ViewUrl  string
	CoverUrl string
	Title    string
}
