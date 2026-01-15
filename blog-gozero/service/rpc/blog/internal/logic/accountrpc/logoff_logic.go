package accountrpclogic

import (
	"context"

	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizcode"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizerr"

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
func (l *LogoffLogic) Logoff(in *accountrpc.LogoffReq) (*accountrpc.LogoffResp, error) {
	// 验证用户是否存在
	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUserNotExist, err.Error())
	}

	err = l.svcCtx.Gorm.Transaction(func(tx *gorm.DB) error {
		err = l.logoff(l.ctx, tx, user.UserId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.OnlineUserService.Logout(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	return &accountrpc.LogoffResp{}, nil
}

func (l *LogoffLogic) logoff(ctx context.Context, tx *gorm.DB, uid string) error {
	// 删除用户账号
	_, err := l.svcCtx.TUserModel.WithTransaction(tx).Deletes(ctx, "user_id = ?", uid)
	if err != nil {
		return err
	}

	// 删除用户角色
	_, err = l.svcCtx.TUserRoleModel.WithTransaction(tx).Deletes(ctx, "user_id = ?", uid)
	if err != nil {
		return err
	}

	return nil
}
