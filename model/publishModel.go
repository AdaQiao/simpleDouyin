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
type UserIdToken struct {
	UserId int64
	Token  string
}
type FilenameAndFilepath struct {
	FileName string
	FilePath string
}
type CoverAndVideoURL struct {
	CoverURL string
	VideoURL string
}
