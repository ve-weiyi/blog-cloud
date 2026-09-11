package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
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
func (l *UpdateChatLogic) UpdateChat(in *discussionrpc.UpdateChatRequest) (*discussionrpc.UpdateChatResponse, error) {
	entity, err := l.svcCtx.TChatModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.Nickname = in.Nickname
	entity.Avatar = in.Avatar
	entity.Type = in.Type
	entity.Content = in.Content
	entity.Status = in.Status

	_, err = l.svcCtx.TChatModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.UpdateChatResponse{Success: true}, nil
}
