package userservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"
)

type ListUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询用户列表
func (l *ListUsersLogic) ListUsers(in *userrpc.ListUsersRequest) (*userrpc.ListUsersResponse, error) {
	var opts []queryx.Option
	if in.ListQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.ListQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.ListQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.ListQuery.Sorts...))
	}

	if in.Username != nil && *in.Username != "" {
		opts = append(opts, queryx.WithCondition("username = ?", *in.Username))
	}
	if in.Email != nil && *in.Email != "" {
		opts = append(opts, queryx.WithCondition("email = ?", *in.Email))
	}
	if in.Mobile != nil && *in.Mobile != "" {
		opts = append(opts, queryx.WithCondition("mobile = ?", *in.Mobile))
	}
	if in.Status != nil {
		opts = append(opts, queryx.WithCondition("status = ?", *in.Status))
	}
	if in.Keyword != nil && *in.Keyword != "" {
		keyword := "%" + *in.Keyword + "%"
		opts = append(opts, queryx.WithCondition("(username LIKE ? OR nickname LIKE ? OR mobile LIKE ? OR email LIKE ?)", keyword, keyword, keyword, keyword))
	}
	if len(in.UserIds) != 0 {
		opts = append(opts, queryx.WithCondition("user_id IN (?)", in.UserIds))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TUserModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*userrpc.User
	for _, v := range records {
		list = append(list, convertTUserToUser(v))
	}

	return &userrpc.ListUsersResponse{
		ListResult: &userrpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
