package accountrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserOnlineListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserOnlineListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserOnlineListLogic {
	return &FindUserOnlineListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查找在线用户列表
func (l *FindUserOnlineListLogic) FindUserOnlineList(in *accountrpc.FindUserListReq) (*accountrpc.FindUserInfoListResp, error) {
	total, err := l.svcCtx.OnlineUserService.GetOnlineUserCount(l.ctx)
	if err != nil {
		return nil, err
	}

	var opts []query.Option
	if in.Paginate != nil {
		opts = append(opts, query.WithPage(int(in.Paginate.Page)))
		opts = append(opts, query.WithSize(int(in.Paginate.PageSize)))
		opts = append(opts, query.WithSorts(in.Paginate.Sorts...))
	}
	page, size, _, _, _ := query.NewQueryBuilder(opts...).Build()
	// 查找在线用户
	uids, err := l.svcCtx.OnlineUserService.GetOnlineUsers(l.ctx, in.Paginate.Page, in.Paginate.PageSize)
	if err != nil {
		return nil, err
	}

	users, err := l.svcCtx.TUserModel.FindALL(l.ctx, "user_id in (?)", uids)
	if err != nil {
		return nil, err
	}

	// 查找用户角色
	urList, err := l.svcCtx.TUserRoleModel.FindALL(l.ctx, "user_id in (?)", uids)
	if err != nil {
		return nil, err
	}

	var ursMap = make(map[string][]int64)
	var roleIds []int64
	for _, item := range urList {
		roleIds = append(roleIds, item.RoleId)
		ursMap[item.UserId] = append(ursMap[item.UserId], item.RoleId)
	}

	// 查找角色信息
	rList, err := l.svcCtx.TRoleModel.FindALL(l.ctx, "id in (?)", roleIds)
	if err != nil {
		return nil, err
	}

	var list []*accountrpc.UserInfo
	for _, item := range users {

		var roles []*model.TRole
		ur, _ := ursMap[item.UserId]
		for _, rid := range ur {
			for _, r := range rList {
				if r.Id == rid {
					roles = append(roles, r)
					break
				}
			}
		}

		list = append(list, convertUserInfoOut(item, roles))
	}

	resp := &accountrpc.FindUserInfoListResp{}
	resp.Pagination = &accountrpc.PageResp{
		Page:     int64(page),
		PageSize: int64(size),
		Total:    total,
	}
	resp.List = list

	return resp, nil
}
