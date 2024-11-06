package talk

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/talkrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTalkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询说说
func NewGetTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTalkLogic {
	return &GetTalkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTalkLogic) GetTalk(req *types.IdReq) (resp *types.TalkBackDTO, err error) {
	in := &talkrpc.IdReq{
		Id: req.Id,
	}
	out, err := l.svcCtx.TalkRpc.GetTalk(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertTalkTypes(out, nil)
	return resp, nil
}
