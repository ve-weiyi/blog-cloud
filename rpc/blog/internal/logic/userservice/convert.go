package userservicelogic

import (
	"context"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

// 数据库模型转换为 Proto 消息

func convertTUserToUser(user *model.TUser) *userrpc.User {
	if user == nil {
		return nil
	}

	return &userrpc.User{
		Id:           user.Id,
		UserId:       user.UserId,
		Username:     user.Username,
		Password:     "",
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Email:        stringPtr(user.Email),
		Mobile:       stringPtr(user.Mobile),
		Status:       user.Status,
		Info:         user.Info,
		RegisterType: user.RegisterType,
		IpAddress:    user.IpAddress,
		IpSource:     user.IpSource,
		CreatedAt:    user.CreatedAt.UnixMilli(),
		UpdatedAt:    user.UpdatedAt.UnixMilli(),
		DeletedAt: func() int64 {
			if user.DeletedAt != nil {
				return user.DeletedAt.UnixMilli()
			}
			return 0
		}(),
	}
}

func stringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func convertTUserToMeInfo(ctx context.Context, svcCtx *svc.ServiceContext, user *model.TUser) *userrpc.MeInfo {
	if user == nil {
		return nil
	}

	var thirdParties []*userrpc.ThirdPartyInfo
	oauthList, _ := svcCtx.TUserOauthModel.FindALL(ctx, "user_id = ?", user.UserId)
	for _, oa := range oauthList {
		thirdParties = append(thirdParties, &userrpc.ThirdPartyInfo{
			Platform:  oa.Platform,
			OpenId:    oa.OpenId,
			Nickname:  oa.Nickname,
			Avatar:    oa.Avatar,
			CreatedAt: oa.CreatedAt.UnixMilli(),
		})
	}

	return &userrpc.MeInfo{
		User:         convertTUserToUser(user),
		ThirdParties: thirdParties,
	}
}
