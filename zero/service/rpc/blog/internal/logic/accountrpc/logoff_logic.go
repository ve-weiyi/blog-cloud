package accountrpclogic

import (
	"context"

	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr/codex"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoffLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoffLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoffLogic {
	return &LogoffLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 注销
func (l *LogoffLogic) Logoff(in *accountrpc.LogoffReq) (*accountrpc.EmptyResp, error) {
	// 验证用户是否存在
	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, apierr.NewApiError(codex.CodeUserNotExist, err.Error())
	}

	err = l.svcCtx.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err = l.logoff(l.ctx, tx, user.UserId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &accountrpc.EmptyResp{}, nil
}

func (l *LogoffLogic) logoff(ctx context.Context, tx *gorm.DB, uid string) (*accountrpc.EmptyResp, error) {
	// 删除用户账号
	_, err := l.svcCtx.TUserModel.WithTransaction(tx).Deletes(ctx, "user_id = ?", uid)
	if err != nil {
		return nil, err
	}

	// 删除用户角色
	_, err = l.svcCtx.TUserRoleModel.WithTransaction(tx).Deletes(ctx, "user_id = ?", uid)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
