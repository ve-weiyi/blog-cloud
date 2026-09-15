package apiutils

import (
	"context"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/guestservice"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"
)

func GetUserInfos(ctx context.Context, svcCtx *svc.ServiceContext, uids []string) (map[string]*types.UserInfoVO, error) {
	if len(uids) == 0 {
		return nil, nil
	}

	out, err := svcCtx.UserService.ListUsers(ctx, &userservice.ListUsersRequest{
		UserIds: uids,
	})
	if err != nil {
		return nil, err
	}

	usm := make(map[string]*types.UserInfoVO)
	for _, v := range out.List {
		usm[v.UserId] = &types.UserInfoVO{
			UserId:   v.UserId,
			Username: v.Username,
			Avatar:   v.Avatar,
			Nickname: v.Nickname,
		}
	}

	return usm, nil
}

func GetGuests(ctx context.Context, svcCtx *svc.ServiceContext, tids []string) (map[string]*types.GuestInfoVO, error) {
	if len(tids) == 0 {
		return nil, nil
	}

	out, err := svcCtx.GuestService.ListGuests(ctx, &guestservice.ListGuestsRequest{
		DeviceIds: tids,
	})
	if err != nil {
		return nil, err
	}

	vsm := make(map[string]*types.GuestInfoVO)
	for _, v := range out.List {
		vsm[v.DeviceId] = &types.GuestInfoVO{
			DeviceId:  v.DeviceId,
			Os:        v.Os,
			Browser:   v.Browser,
			IpAddress: v.IpAddress,
			IpSource:  v.IpSource,
		}
	}

	return vsm, nil
}
