package accountrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改用户状态
func (l *UpdateUserStatusLogic) UpdateUserStatus(in *accountrpc.UpdateUserStatusReq) (*accountrpc.EmptyResp, error) {
	ua, err := l.svcCtx.TUserModel.First(l.ctx, "id = ?", in.UserId)
	if err != nil {
		return nil, err
	}

	ua.Status = in.Status

	_, err = l.svcCtx.TUserModel.Save(l.ctx, ua)
	if err != nil {
		return nil, err
	}

	return &accountrpc.EmptyResp{}, nil
}
