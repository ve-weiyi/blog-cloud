package talk

import (
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/discussionservice"
)

// convertTalkVO 把说说实体映射为管理端视图。
// 作者信息由调用方批量补齐后传入，避免逐条查询。
func convertTalkVO(in *discussionservice.Talk, usm map[string]*types.UserInfoVO) *types.TalkVO {
	if in == nil {
		return nil
	}

	return &types.TalkVO{
		Id:           in.Id,
		UserId:       in.UserId,
		Content:      in.Content,
		ImgList:      in.Images,
		IsTop:        in.IsTop,
		Status:       in.Status,
		LikeCount:    in.LikeCount,
		CommentCount: in.CommentCount,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		UserInfo:     usm[in.UserId],
	}
}
