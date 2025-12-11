package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateChatLogic {
	return &UpdateChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新聊天记录
func (l *UpdateChatLogic) UpdateChat(in *messagerpc.UpdateChatReq) (*messagerpc.ChatDetailsResp, error) {
	entity, err := l.svcCtx.TChatModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.Type = in.Type
	entity.Content = in.Content
	entity.Status = in.Status

	_, err = l.svcCtx.TChatModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertChatOut(entity), nil
}
