package types

type QueryTalkReq struct {
	PageQuery
}

type Talk struct {
	Id           int64       `json:"id"`            // 说说ID
	UserId       string      `json:"user_id"`       // 用户ID
	Content      string      `json:"content"`       // 评论内容
	ImgList      []string    `json:"img_list"`      // 图片URL列表
	IsTop        int64       `json:"is_top"`        // 是否置顶
	Status       int64       `json:"status"`        // 状态 1.公开 2.私密
	LikeCount    int64       `json:"like_count"`    // 点赞量
	CommentCount int64       `json:"comment_count"` // 评论量
	CreatedAt    int64       `json:"created_at"`    // 创建时间
	UpdatedAt    int64       `json:"updated_at"`    // 更新时间
	UserInfo     *UserInfoVO `json:"user_info"`     // 用户信息
}
