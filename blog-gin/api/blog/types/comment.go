package types

type NewCommentReq struct {
	TopicId        int64  `json:"topic_id,optional"`      // 主题id
	ParentId       int64  `json:"parent_id,optional"`     // 父评论id
	ReplyId        int64  `json:"reply_id,optional"`      // 会话id
	ReplyUserId    string `json:"reply_user_id,optional"` // 回复用户id
	CommentContent string `json:"comment_content"`        // 评论内容
	Status         int64  `json:"status,optional"`        // 状态
	Type           int64  `json:"type"`                   // 评论类型 1.文章 2.友链 3.说说
}

type QueryCommentReq struct {
	PageQuery
	TopicId  int64 `json:"topic_id,optional"`  // 主题id
	ParentId int64 `json:"parent_id,optional"` // 父评论id
	Type     int64 `json:"type,optional"`      // 评论类型 1.文章 2.友链 3.说说
}

type UpdateCommentReq struct {
	Id             int64  `json:"id"`                     // 主键
	ReplyUserId    string `json:"reply_user_id,optional"` // 回复用户id
	CommentContent string `json:"comment_content"`        // 评论内容
	Status         int64  `json:"status,optional"`        // 状态
}
