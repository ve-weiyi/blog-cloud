package message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
)

type CreateMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建留言
func NewCreateMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMessageLogic {
	return &CreateMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMessageLogic) CreateMessage(req *types.CreateMessageReq) (resp *types.EmptyResp, err error) {
	uid, _ := metax.GetApiUserIdFromCtx(l.ctx)
	did, _ := metax.GetApiDeviceIdFromCtx(l.ctx)

	_, err = l.svcCtx.DiscussionService.CreateMessage(l.ctx, &discussionservice.CreateMessageRequest{
		UserId:         uid,
		DeviceId:       did,
		MessageContent: req.MessageContent,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
