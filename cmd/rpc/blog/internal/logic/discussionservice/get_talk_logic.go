package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
)

type GetTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTalkLogic {
	return &GetTalkLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetTalkLogic) GetTalk(in *discussionrpc.GetTalkRequest) (*discussionrpc.GetTalkResponse, error) {
	conditions := "`id` = ?"
	args := []any{in.Id}
	if in.PublicOnly {
		conditions += " AND `status` = ?"
		args = append(args, enums.TalkStatusPublic)
	}
	entity, err := l.svcCtx.TTalkModel.FindOne(l.ctx, "", conditions, args...)
	if err != nil {
		return nil, err
	}
	lm, lcm, err := likeStats(l.ctx, l.svcCtx, enums.LikeTargetTalk, []int64{entity.Id})
	if err != nil {
		return nil, err
	}

	return &discussionrpc.GetTalkResponse{Talk: convertTalkOut(entity, lm, lcm)}, nil
}
