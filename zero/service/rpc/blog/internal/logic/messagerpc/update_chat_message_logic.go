package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/messagerpc"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateChatMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateChatMessageLogic {
	return &UpdateChatMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新聊天记录
func (l *UpdateChatMessageLogic) UpdateChatMessage(in *messagerpc.ChatMessageNewReq) (*messagerpc.ChatMessageDetails, error) {
	entity := convertChatMessageIn(in)

	_, err := l.svcCtx.TChatMessageModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertChatMessageOut(entity), nil
}
