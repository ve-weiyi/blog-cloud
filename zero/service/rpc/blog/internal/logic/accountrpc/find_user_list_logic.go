package accountrpclogic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"
)

type FindUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserListLogic {
	return &FindUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查找用户列表
func (l *FindUserListLogic) FindUserList(in *accountrpc.FindUserListReq) (*accountrpc.FindUserListResp, error) {
	page, size, sorts, conditions, params := convertQuery(in)

	result, err := l.svcCtx.TUserModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	total, err := l.svcCtx.TUserModel.FindCount(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*accountrpc.User
	for _, item := range result {
		list = append(list, convertUserOut(item))
	}

	resp := &accountrpc.FindUserListResp{}
	resp.Total = total
	resp.List = list

	return resp, nil
}

func convertQuery(in *accountrpc.FindUserListReq) (page int, size int, sorts string, conditions string, params []interface{}) {
	page = int(in.Page)
	size = int(in.PageSize)
	sorts = strings.Join(in.Sorts, ",")
	if sorts == "" {
		sorts = "id desc"
	}

	if in.Username != "" {
		conditions += "username like ?"
		params = append(params, "%"+in.Username+"%")
	}

	if in.Nickname != "" {
		if conditions != "" {
			conditions += " and "
		}
		conditions += "nickname like ?"
		params = append(params, "%"+in.Nickname+"%")
	}

	if in.Email != "" {
		if conditions != "" {
			conditions += " and "
		}
		conditions += "email like ?"
		params = append(params, "%"+in.Email+"%")
	}

	if in.Status != 0 {
		if conditions != "" {
			conditions += " and "
		}
		conditions += "status = ?"
		params = append(params, in.Status)
	}

	if len(in.UserIds) != 0 {
		if conditions != "" {
			conditions += " and "
		}
		conditions += "id in (?)"
		params = append(params, in.UserIds)
	}

	return page, size, sorts, conditions, params
}

func convertUserOut(in *model.TUser) (out *accountrpc.User) {

	out = &accountrpc.User{
		UserId:    in.UserId,
		Username:  in.Username,
		Nickname:  in.Nickname,
		Avatar:    in.Avatar,
		Email:     in.Email,
		Phone:     in.Phone,
		Info:      in.Info,
		Status:    in.Status,
		LoginType: in.LoginType,
		IpAddress: in.IpAddress,
		IpSource:  in.IpSource,
		CreatedAt: in.CreatedAt.Unix(),
		UpdatedAt: in.UpdatedAt.Unix(),
	}

	return out
}
