package userservicelogic

import (
	"context"
	"errors"
	"fmt"
	"html"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
	"github.com/ve-weiyi/blog-cloud/infra/store/otpx"
	"github.com/ve-weiyi/vkit/adapter/emailx"
	"github.com/ve-weiyi/vkit/x/validatex"
)

type BindMeEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindMeEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMeEmailLogic {
	return &BindMeEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 绑定当前用户邮箱
func (l *BindMeEmailLogic) BindMeEmail(in *userrpc.BindMeEmailRequest) (*userrpc.BindMeEmailResponse, error) {
	userID, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil || userID == "" {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "缺少用户身份")
	}
	if !validatex.IsEmail(in.Email) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "邮箱格式不正确")
	}

	if err := l.verifyCode(in.Email, in.VerifyCode); err != nil {
		return nil, err
	}

	existing, err := l.svcCtx.TUserModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil && existing.UserId != userID {
		return nil, bizerr.NewBizError(bizcode.CodeResourceAlreadyExist, "邮箱已被绑定")
	}

	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, userID)
	if err != nil {
		return nil, err
	}
	// 改之前先记下旧邮箱：告警必须送到它，否则会话被盗时受害者什么都收不到
	oldEmail := ""
	if user.Email != nil {
		oldEmail = *user.Email
	}

	err = l.svcCtx.GormDB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := l.svcCtx.TUserModel.WithTx(tx).UpdateFields(l.ctx, map[string]interface{}{
			"email": in.Email,
		}, "user_id = ?", userID); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return bizerr.NewBizError(bizcode.CodeResourceAlreadyExist, "邮箱已被绑定")
			}
			return err
		}
		return l.svcCtx.Notify.Emit(l.ctx, tx, notification.Request{
			Scene:       notify.SceneAccountSecurity,
			ObjectType:  notify.ObjectUser,
			ObjectID:    user.Id,
			Audience:    notification.AudienceUserIDs,
			AudienceIDs: []string{userID},
			Payload:     notify.Payload{Action: notify.SecurityEmailChanged},
		})
	})
	if err != nil {
		return nil, err
	}

	l.sendEmailChangedAlert(oldEmail, in.Email)

	return &userrpc.BindMeEmailResponse{Success: true}, nil
}

// emailAlertTimeout 旧邮箱告警的等待上限。
const emailAlertTimeout = 5 * time.Second

// sendEmailChangedAlert 直发一封告警到**旧邮箱**。
//
// 不能走通知管线：管线在扇出时才从 t_user 解析地址，那时邮箱已经变成新的，
// 受害者的旧邮箱永远收不到——而这正是被盗号改邮箱时要告警的那一方。
// 直发失败只记日志，不影响改邮箱本身。
func (l *BindMeEmailLogic) sendEmailChangedAlert(oldEmail, newEmail string) {
	if oldEmail == "" || oldEmail == newEmail || l.svcCtx.EmailDeliver == nil {
		return
	}

	content := fmt.Sprintf(
		`<p>你的账号邮箱已由 %s 变更为 %s。</p><p>如果这不是你本人的操作，请立即修改密码并联系管理员。</p>`,
		html.EscapeString(oldEmail), html.EscapeString(newEmail))
	// 直发是同步的（安全兜底不能丢），但要给等待设上限：SMTP 卡住时不能连带拖住改邮箱接口。
	// 超时只是不再等待，发送本身仍在进行
	done := make(chan error, 1)
	go func() {
		done <- l.svcCtx.EmailDeliver.DeliveryEmail(&emailx.EmailMessage{
			To:      []string{oldEmail},
			CC:      []string{},
			Subject: "账号邮箱已变更",
			Content: content,
		})
	}()

	select {
	case err := <-done:
		if err != nil {
			l.Errorf("旧邮箱告警发送失败: %v", err)
		}
	case <-time.After(emailAlertTimeout):
		l.Errorf("旧邮箱告警发送超时（%s），已放弃等待", emailAlertTimeout)
	}
}

func (l *BindMeEmailLogic) verifyCode(email, code string) error {
	key := fmt.Sprintf("email:code:%s:%s", enums.CodeSceneBind, email)
	ok, err := l.svcCtx.OTPStore.Verify(l.ctx, key, code)
	if errors.Is(err, otpx.ErrTooManyAttempts) {
		return bizerr.NewBizError(bizcode.CodeVerifyCodeError, "验证码错误次数过多，请重新获取")
	}
	if err != nil {
		return err
	}
	if !ok {
		return bizerr.NewBizError(bizcode.CodeVerifyCodeError, "验证码错误或已过期")
	}
	return nil
}
