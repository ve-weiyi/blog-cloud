package types

type Comment struct {
	Id               int64           `json:"id"`                 // 评论id
	TopicId          int64           `json:"topic_id"`           // 主题id
	ParentId         int64           `json:"parent_id"`          // 父评论id
	ReplyMsgId       int64           `json:"reply_msg_id"`       // 会话id
	UserId           string          `json:"user_id"`            // 用户id
	ReplyUserId      string          `json:"reply_user_id"`      // 被回复用户id
	CommentContent   string          `json:"comment_content"`    // 评论内容
	Type             int64           `json:"type"`               // 评论类型 1.文章 2.友链 3.说说
	CreatedAt        int64           `json:"created_at"`         // 评论时间
	LikeCount        int64           `json:"like_count"`         // 点赞数
	User             *UserInfoVO     `json:"user"`               // 评论用户
	ReplyUser        *UserInfoVO     `json:"reply_user"`         // 被回复评论用户
	ReplyCount       int64           `json:"reply_count"`        // 回复量
	CommentReplyList []*CommentReply `json:"comment_reply_list"` // 评论回复列表
}

type CommentNewReq struct {
	TopicId        int64  `json:"topic_id,optional"`      // 主题id
	ParentId       int64  `json:"parent_id,optional"`     // 父评论id
	ReplyMsgId     int64  `json:"reply_msg_id,optional"`  // 会话id
	ReplyUserId    string `json:"reply_user_id,optional"` // 回复用户id
	CommentContent string `json:"comment_content"`        // 评论内容
	Type           int64  `json:"type"`                   // 评论类型 1.文章 2.友链 3.说说
	Status         int64  `json:"status,optional"`        // 状态 0.正常 1.已编辑 2.已删除
}

type CommentQueryReq struct {
	PageQuery
	TopicId  int64 `json:"topic_id,optional"`  // 主题id
	ParentId int64 `json:"parent_id,optional"` // 父评论id
	Type     int64 `json:"type,optional"`      // 评论类型 1.文章 2.友链 3.说说
}

type UpdateCommentReq struct {
	Id             int64  `json:"id"`                     // 主键
	ReplyUserId    string `json:"reply_user_id,optional"` // 回复用户id
	CommentContent string `json:"comment_content"`        // 评论内容
	Status         int64  `json:"status,optional"`        // 状态 0.正常 1.已编辑 2.已删除
}
