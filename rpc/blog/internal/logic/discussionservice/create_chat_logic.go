package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type CreateChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateChatLogic {
	return &CreateChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建聊天记录
func (l *CreateChatLogic) CreateChat(in *discussionrpc.CreateChatRequest) (*discussionrpc.CreateChatResponse, error) {
	entity := &model.TChat{
		UserId:    in.UserId,
		DeviceId:  in.DeviceId,
		IpAddress: in.IpAddress,
		IpSource:  in.IpSource,
		Nickname:  in.Nickname,
		Avatar:    in.Avatar,
		Type:      in.Type,
		Content:   in.Content,
		Status:    in.Status,
	}
	_, err := l.svcCtx.TChatModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.CreateChatResponse{Id: entity.Id}, nil
}
