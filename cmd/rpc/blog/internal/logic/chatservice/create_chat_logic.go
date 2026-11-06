package chatservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/chatrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/chat"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/model"
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
func (l *CreateChatLogic) CreateChat(in *chatrpc.CreateChatRequest) (*chatrpc.CreateChatResponse, error) {
	if !chat.ValidType(in.Type) {
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "不支持的消息类型")
	}
	content, err := chat.NormalizeContent(in.Content)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "消息内容不合法")
	}

	entity := &model.TChat{
		UserId:    in.UserId,
		DeviceId:  in.DeviceId,
		IpAddress: in.IpAddress,
		IpSource:  in.IpSource,
		Nickname:  in.Nickname,
		Avatar:    in.Avatar,
		Type:      in.Type,
		Content:   content,
		Status:    in.Status,
	}
	_, err = l.svcCtx.TChatModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &chatrpc.CreateChatResponse{Id: entity.Id}, nil
}
