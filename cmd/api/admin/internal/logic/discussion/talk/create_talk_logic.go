package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/enrich"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/discussionservice"
)

type CreateTalkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建说说
func NewCreateTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTalkLogic {
	return &CreateTalkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTalkLogic) CreateTalk(req *types.CreateTalkReq) (resp *types.TalkVO, err error) {
	out, err := l.svcCtx.DiscussionService.CreateTalk(l.ctx, &discussionservice.CreateTalkRequest{
		Content: req.Content,
		Images:  req.ImgList,
		IsTop:   req.IsTop,
		Status:  req.Status,
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
