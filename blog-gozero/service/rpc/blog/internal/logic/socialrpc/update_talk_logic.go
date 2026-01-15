package socialrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/socialrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTalkLogic {
	return &UpdateTalkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新说说
func (l *UpdateTalkLogic) UpdateTalk(in *socialrpc.UpdateTalkReq) (*socialrpc.UpdateTalkResp, error) {
	entity, err := l.svcCtx.TTalkModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.UserId = in.UserId
	entity.Content = in.Content
	entity.Images = jsonconv.AnyToJsonNE(in.ImgList)
	entity.IsTop = in.IsTop
	entity.Status = in.Status

	_, err = l.svcCtx.TTalkModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &socialrpc.UpdateTalkResp{
		Talk: convertTalkOut(entity),
	}, nil
}
