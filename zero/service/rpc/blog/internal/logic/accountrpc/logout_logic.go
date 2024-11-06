package accountrpclogic

import (
	"context"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 登出
func (l *LogoutLogic) Logout(in *accountrpc.LogoutReq) (*accountrpc.LogoutResp, error) {
	list, err := l.svcCtx.TUserLoginHistoryModel.FindList(l.ctx, 1, 1, "id desc", "user_id = ?", in.UserId)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return &accountrpc.LogoutResp{
			UserId: in.UserId,
		}, nil
	}

	find := list[0]

	find.LogoutAt = time.Now()
	// 修改登出时间
	_, err = l.svcCtx.TUserLoginHistoryModel.Save(l.ctx, find)
	if err != nil {
		return nil, err
	}

	return &accountrpc.LogoutResp{
		UserId:   in.UserId,
		LogoutAt: find.LogoutAt.Unix(),
	}, nil
}
