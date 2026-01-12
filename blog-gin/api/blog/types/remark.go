package types

type NewRemarkReq struct {
	MessageContent string `json:"message_content"` // 留言内容
}

type QueryRemarkReq struct {
	PageQuery
}
