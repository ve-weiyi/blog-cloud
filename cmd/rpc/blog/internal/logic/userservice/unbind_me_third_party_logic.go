package userservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type UnbindMeThirdPartyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnbindMeThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindMeThirdPartyLogic {
	return &UnbindMeThirdPartyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 解绑当前用户第三方平台
func (l *UnbindMeThirdPartyLogic) UnbindMeThirdParty(in *userrpc.UnbindMeThirdPartyRequest) (*userrpc.UnbindMeThirdPartyResponse, error) {
	userId, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	oa, err := l.svcCtx.TUserOauthModel.FindOneByUserIdPlatform(l.ctx, userId, in.Platform)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, fmt.Sprintf("未找到平台 %s 的绑定记录", in.Platform))
	}

	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, userId)
	if err != nil {
		return nil, err
	}
	others, err := l.svcCtx.TUserOauthModel.FindAll(l.ctx, "user_id = ? AND platform != ?", userId, in.Platform)
	if err != nil {
		return nil, err
	}

	hasPassword := strings.TrimSpace(user.Password) != ""
	hasEmail := user.Email != nil && *user.Email != ""
	hasMobile := user.Mobile != nil && *user.Mobile != ""
	if !hasPassword && !hasEmail && !hasMobile && len(others) == 0 {
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "至少保留一种登录方式")
	}

	_, err = l.svcCtx.TUserOauthModel.Delete(l.ctx, oa.Id)
	if err != nil {
		return nil, err
	}

	return &userrpc.UnbindMeThirdPartyResponse{
		Success: true,
	}, nil
}
