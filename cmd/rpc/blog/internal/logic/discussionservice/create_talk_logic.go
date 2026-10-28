package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/model"
	"github.com/ve-weiyi/vkit/x/jsonv"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
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
	images := jsonv.AnyToJsonNE(in.Images)
	entity := &model.TTalk{
		UserId:  in.UserId,
		Content: in.Content,
		Images:  &images,
		IsTop:   in.IsTop,
		Status:  in.Status,
	}
	_, err := l.svcCtx.TTalkModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	lm, lcm, err := likeStats(l.ctx, l.svcCtx, enums.LikeTargetTalk, []int64{entity.Id})
	if err != nil {
		return nil, err
	}

	return &discussionrpc.CreateTalkResponse{Id: entity.Id, Talk: convertTalkOut(entity, lm, lcm)}, nil
}
