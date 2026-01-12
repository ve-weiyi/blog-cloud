package messagerpclogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
)

type AddChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddChatLogic {
	return &AddChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建聊天记录
func (l *AddChatLogic) AddChat(in *messagerpc.AddChatReq) (*messagerpc.ChatDetailsResp, error) {
	entity := &model.TChat{
		Id:         0,
		UserId:     in.UserId,
		TerminalId: in.TerminalId,
		Nickname:   in.Nickname,
		Avatar:     in.Avatar,
		IpAddress:  in.IpAddress,
		IpSource:   in.IpSource,
		Type:       in.Type,
		Content:    in.Content,
		Status:     in.Status,
		CreatedAt:  time.Now(),
	}

	_, err := l.svcCtx.TChatModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertChatOut(entity), nil
}

func convertChatOut(in *model.TChat) (out *messagerpc.ChatDetailsResp) {
	out = &messagerpc.ChatDetailsResp{
		Id:         in.Id,
		UserId:     in.UserId,
		TerminalId: in.TerminalId,
		Nickname:   in.Nickname,
		Avatar:     in.Avatar,
		IpAddress:  in.IpAddress,
		IpSource:   in.IpSource,
		Type:       in.Type,
		Content:    in.Content,
		Status:     in.Status,
		CreatedAt:  in.CreatedAt.Unix(),
		UpdatedAt:  in.UpdatedAt.Unix(),
	}

	return out
}
