package userservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/auth/passwordx"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type ResetUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	return &ResetUserPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 重置用户密码
func (l *ResetUserPasswordLogic) ResetUserPassword(in *userrpc.ResetUserPasswordRequest) (*userrpc.ResetUserPasswordResponse, error) {
	if err := passwordx.ValidatePassword(in.NewPassword); err != nil {
		return nil, err
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	// 操作人要记进安全通知：取不到身份说明网关没注入主体，继续写空值会发出无主告警
	operator, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "未登录")
	}

	// 管理员重置他人密码是安全事件，必须通知本人；触发者是管理员，不排除收件人
	err = l.svcCtx.GormDB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := l.svcCtx.TUserModel.WithTx(tx).UpdateFields(l.ctx, map[string]interface{}{
			"password": string(hashedPassword),
		}, "user_id = ?", in.UserId); err != nil {
			return err
		}
		return l.svcCtx.Notify.Emit(l.ctx, tx, notification.Request{
			Scene:       notify.SceneAccountSecurity,
			ActorID:     operator,
			ObjectType:  notify.ObjectUser,
			ObjectID:    user.Id,
			Audience:    notification.AudienceUserIDs,
			AudienceIDs: []string{in.UserId},
			Payload:     notify.Payload{Action: notify.SecurityPasswordChanged},
		})
	})
	if err != nil {
		return nil, err
	}

	return &userrpc.ResetUserPasswordResponse{
		Success: true,
	}, nil
}
