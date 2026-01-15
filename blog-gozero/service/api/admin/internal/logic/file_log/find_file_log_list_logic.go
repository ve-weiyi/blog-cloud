package file_log

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/common/apiutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindFileLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询文件日志
func NewFindFileLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindFileLogListLogic {
	return &FindFileLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindFileLogListLogic) FindFileLogList(req *types.QueryFileLogReq) (resp *types.PageResp, err error) {
	in := &syslogrpc.FindFileLogListReq{
		Paginate: &syslogrpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		FilePath: req.FilePath,
		FileName: req.FileName,
		FileType: req.FileType,
	}

	out, err := l.svcCtx.SyslogRpc.FindFileLogList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	usm, err := apiutils.BatchQuery(out.List,
		func(v *syslogrpc.FileLog) string {
			return v.UserId
		},
		func(ids []string) (map[string]*types.UserInfoVO, error) {
			return apiutils.GetUserInfos(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}

	// 查询访客信息
	vsm, err := apiutils.BatchQuery(out.List,
		func(v *syslogrpc.FileLog) string {
			return v.TerminalId
		},
		func(ids []string) (map[string]*types.ClientInfoVO, error) {
			return apiutils.GetVisitorInfos(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}

	var list []*types.FileLogBackVO
	for _, v := range out.List {
		list = append(list, &types.FileLogBackVO{
			Id:         v.Id,
			UserId:     v.UserId,
			FilePath:   v.FilePath,
			FileName:   v.FileName,
			FileType:   v.FileType,
			FileSize:   v.FileSize,
			FileMd5:    v.FileMd5,
			FileUrl:    v.FileUrl,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
			UserInfo:   usm[v.UserId],
			ClientInfo: vsm[v.TerminalId],
		})
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
