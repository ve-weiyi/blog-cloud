package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/socialservice"
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
	out, err := l.svcCtx.SocialService.CreateTalk(l.ctx, &socialservice.CreateTalkRequest{
		Content: req.Content,
		Images:  req.ImgList,
		IsTop:   req.IsTop,
		Status:  req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &types.TalkVO{
		Id:      out.Id,
		Content: req.Content,
		ImgList: req.ImgList,
		IsTop:   req.IsTop,
		Status:  req.Status,
	}, nil
}
