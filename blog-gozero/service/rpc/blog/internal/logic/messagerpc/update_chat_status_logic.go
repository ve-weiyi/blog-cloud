package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateChatStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateChatStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateChatStatusLogic {
	return &UpdateChatStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新聊天记录状态
func (l *UpdateChatStatusLogic) UpdateChatStatus(in *messagerpc.UpdateChatStatusReq) (*messagerpc.ChatDetailsResp, error) {
	entity, err := l.svcCtx.TChatModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.Status = in.Status

	_, err = l.svcCtx.TChatModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertChatOut(entity), nil
}
