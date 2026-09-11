package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/socialservice"
)

type LikeTalkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 点赞说说
func NewLikeTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeTalkLogic {
	return &LikeTalkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeTalkLogic) LikeTalk(req *types.LikeTalkReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.SocialService.LikeTalk(l.ctx, &socialservice.LikeTalkRequest{
		Id: req.TalkId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.EmptyResp{}
	return
}
