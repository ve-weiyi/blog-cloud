package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/x/jsonv"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
)

type UpdateTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTalkLogic {
	return &UpdateTalkLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateTalkLogic) UpdateTalk(in *discussionrpc.UpdateTalkRequest) (*discussionrpc.UpdateTalkResponse, error) {
	entity, err := l.svcCtx.TTalkModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	entity.Content = in.Content
	entity.Images = jsonv.AnyToJsonNE(in.Images)
	entity.IsTop = in.IsTop
	entity.Status = in.Status

	_, err = l.svcCtx.TTalkModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}
	return &discussionrpc.UpdateTalkResponse{Success: true}, nil
}
