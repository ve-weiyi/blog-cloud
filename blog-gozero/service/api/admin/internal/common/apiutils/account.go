package apiutils

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

func GetUserInfos(ctx context.Context, svcCtx *svc.ServiceContext, uids []string) (map[string]*types.UserInfoVO, error) {
	if len(uids) == 0 {
		return nil, nil
	}
	users, err := svcCtx.AccountRpc.FindUserList(ctx, &accountrpc.FindUserListReq{
		UserIds: uids,
	})
	if err != nil {
		return nil, err
	}

	usm := make(map[string]*types.UserInfoVO)
	for _, v := range users.List {
		var ext types.UserInfoExt
		if v.Info != "" {
			err = jsonconv.JsonToAny(v.Info, &ext)
			if err != nil {
				return nil, err
			}
		}

		usm[v.UserId] = &types.UserInfoVO{
			UserId:      v.UserId,
			Username:    v.Username,
			Avatar:      v.Avatar,
			Nickname:    v.Nickname,
			UserInfoExt: ext,
		}
	}

	return usm, err
}
