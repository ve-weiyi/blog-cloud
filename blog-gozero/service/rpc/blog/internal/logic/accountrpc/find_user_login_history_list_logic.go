package accountrpclogic

import (
	"context"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLoginHistoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLoginHistoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLoginHistoryListLogic {
	return &FindUserLoginHistoryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户登录历史
func (l *FindUserLoginHistoryListLogic) FindUserLoginHistoryList(in *accountrpc.FindLoginHistoryListReq) (*accountrpc.FindLoginHistoryListResp, error) {
	page, size, sorts, conditions, params := convertUserLoginHistoryQuery(in)

	result, err := l.svcCtx.TUserLoginHistoryModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*accountrpc.UserLoginHistory
	for _, item := range result {
		list = append(list, convertUserLoginHistoryOut(item))
	}

	total, err := l.svcCtx.TUserLoginHistoryModel.FindCount(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	resp := &accountrpc.FindLoginHistoryListResp{}
	resp.List = list
	resp.Total = total
	return resp, nil
}

func convertUserLoginHistoryQuery(in *accountrpc.FindLoginHistoryListReq) (page int, size int, sorts string, conditions string, params []interface{}) {
	page = int(in.Page)
	size = int(in.PageSize)
	sorts = strings.Join(in.Sorts, ",")
	if sorts == "" {
		sorts = "id desc"
	}
	if in.UserId != "" {
		conditions += "user_id = ?"
		params = append(params, in.UserId)
	}

	return page, size, sorts, conditions, params
}

func convertUserLoginHistoryOut(in *model.TUserLoginHistory) (out *accountrpc.UserLoginHistory) {
	out = &accountrpc.UserLoginHistory{
		Id:        in.Id,
		UserId:    in.UserId,
		LoginType: in.LoginType,
		Agent:     in.Agent,
		IpAddress: in.IpAddress,
		IpSource:  in.IpSource,
		LoginAt:   in.LoginAt.Unix(),
		LogoutAt:  in.LogoutAt.Unix(),
	}

	return out
}
