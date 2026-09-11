package user

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/guestservice"
)

type GetOnlineUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取在线用户列表
func NewGetOnlineUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOnlineUsersLogic {
	return &GetOnlineUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOnlineUsersLogic) GetOnlineUsers(req *types.EmptyReq) (resp *types.OnlineUserListResp, err error) {
	// 从 blog:online:user 查询博客前台在线用户
	rawUsers, err := l.queryOnlineUsers(l.ctx)
	if err != nil {
		return nil, err
	}

	if len(rawUsers) == 0 {
		return &types.OnlineUserListResp{List: []*types.OnlineUserItem{}}, nil
	}

	// 分离注册用户(UUID)和游客(deviceId)
	userIds := make([]string, 0)
	deviceIds := make([]string, 0)
	userMap := make(map[string]int64) // userId/deviceId -> lastActiveAt
	for _, u := range rawUsers {
		userMap[u.UserId] = u.LastActiveAt
		if isUUID(u.UserId) {
			userIds = append(userIds, u.UserId)
		} else {
			deviceIds = append(deviceIds, u.UserId)
		}
	}

	// 批量查用户信息
	userInfoMap := make(map[string]*userservice.User)
	if len(userIds) > 0 {
		out, err := l.svcCtx.UserService.ListUsers(l.ctx, &userservice.ListUsersRequest{
			UserIds: userIds,
		})
		if err == nil && out != nil {
			for _, u := range out.List {
				userInfoMap[u.UserId] = u
			}
		}
	}

	// 批量查游客信息
	clientInfoMap := make(map[string]*guestservice.Guest)
	if len(deviceIds) > 0 {
		out, err := l.svcCtx.GuestService.ListGuests(l.ctx, &guestservice.ListGuestsRequest{
			DeviceIds: deviceIds,
		})
		if err == nil && out != nil {
			for _, c := range out.List {
				clientInfoMap[c.DeviceId] = c
			}
		}
	}

	list := make([]*types.OnlineUserItem, 0)
	for _, u := range rawUsers {
		item := &types.OnlineUserItem{
			LastActiveAt: u.LastActiveAt,
		}

		if ui, ok := userInfoMap[u.UserId]; ok {
			item.UserInfo = &types.UserInfoVO{
				UserId:   ui.UserId,
				Username: ui.Username,
				Nickname: ui.Nickname,
				Avatar:   ui.Avatar,
				UserType: "user",
			}
		}

		if ci, ok := clientInfoMap[u.UserId]; ok {
			item.GuestInfo = &types.GuestInfoVO{
				DeviceId:  ci.DeviceId,
				Os:        ci.Os,
				Browser:   ci.Browser,
				IpAddress: ci.IpAddress,
				IpSource:  ci.IpSource,
			}
		}

		list = append(list, item)
	}

	return &types.OnlineUserListResp{List: list}, nil
}

// queryOnlineUsers 直接从 Redis ZSET 查询在线用户
func (l *GetOnlineUsersLogic) queryOnlineUsers(ctx context.Context) ([]*rawOnlineUser, error) {
	key := cachekey.OnlineUserKey
	threshold := time.Now().Add(-60 * time.Second).UnixMilli()

	results, err := l.svcCtx.RedisClient.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min: strconv.FormatInt(threshold, 10),
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, err
	}

	users := make([]*rawOnlineUser, 0, len(results))
	for _, z := range results {
		member, ok := z.Member.(string)
		if !ok {
			continue
		}
		users = append(users, &rawOnlineUser{
			UserId:       member,
			LastActiveAt: int64(z.Score),
		})
	}
	return users, nil
}

type rawOnlineUser struct {
	UserId       string
	LastActiveAt int64
}

// isUUID 简单判断是否 UUID 格式 (8-4-4-4-12)
func isUUID(s string) bool {
	return len(s) == 36 && strings.Count(s, "-") == 4
}
