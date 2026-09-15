package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/x/jsonconv"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
)

type CreateTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTalkLogic {
	return &CreateTalkLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateTalkLogic) CreateTalk(in *discussionrpc.CreateTalkRequest) (*discussionrpc.CreateTalkResponse, error) {
	entity := &model.TTalk{
		UserId:  in.UserId,
		Content: in.Content,
		Images:  jsonconv.AnyToJsonNE(in.Images),
		IsTop:   in.IsTop,
		Status:  in.Status,
	}
	_, err := l.svcCtx.TTalkModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}
	return &discussionrpc.CreateTalkResponse{Id: entity.Id}, nil
}
