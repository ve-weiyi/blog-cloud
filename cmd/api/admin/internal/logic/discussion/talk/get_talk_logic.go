package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/enrich"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/discussionservice"
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

func (l *GetTalkLogic) GetTalk(req *types.GetTalkReq) (resp *types.TalkVO, err error) {
	out, err := l.svcCtx.DiscussionService.GetTalk(l.ctx, &discussionservice.GetTalkRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	usm, err := enrich.GetUserInfos(l.ctx, l.svcCtx, []string{out.Talk.UserId})
	if err != nil {
		return nil, err
	}

	return convertTalkVO(out.Talk, usm), nil
}
