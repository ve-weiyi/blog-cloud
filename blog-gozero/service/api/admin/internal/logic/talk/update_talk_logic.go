package talk

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/talkrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTalkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新说说
func NewUpdateTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTalkLogic {
	return &UpdateTalkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTalkLogic) UpdateTalk(req *types.TalkNewReq) (resp *types.TalkBackVO, err error) {
	in := &talkrpc.TalkNewReq{
		Id:      req.Id,
		UserId:  cast.ToString(l.ctx.Value(restx.HeaderUid)),
		Content: req.Content,
		ImgList: req.ImgList,
		IsTop:   req.IsTop,
		Status:  req.Status,
	}

	out, err := l.svcCtx.TalkRpc.UpdateTalk(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertTalkTypes(out, nil)
	return resp, nil
}
