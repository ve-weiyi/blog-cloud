package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/common/apiutils"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/socialservice"
)

type GetTalkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取说说详情
func NewGetTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTalkLogic {
	return &GetTalkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTalkLogic) GetTalk(req *types.GetTalkReq) (resp *types.Talk, err error) {
	out, err := l.svcCtx.SocialService.GetTalk(l.ctx, &socialservice.GetTalkRequest{
		Id: req.TalkId,
	})
	if err != nil {
		return nil, err
	}

	usm, err := apiutils.GetUserInfos(l.ctx, l.svcCtx, []string{out.Talk.UserId})
	if err != nil {
		return nil, err
	}

	resp = &types.Talk{
		Id:           out.Talk.Id,
		UserId:       out.Talk.UserId,
		Content:      out.Talk.Content,
		ImgList:      out.Talk.Images,
		IsTop:        out.Talk.IsTop,
		Status:       out.Talk.Status,
		LikeCount:    out.Talk.LikeCount,
		CommentCount: out.Talk.CommentCount,
		CreatedAt:    out.Talk.CreatedAt,
		UpdatedAt:    out.Talk.UpdatedAt,
		UserInfo:     usm[out.Talk.UserId],
	}
	return
}
