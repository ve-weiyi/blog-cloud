package websiterpclogic

import (
	"context"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/websiterpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserDailyVisitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserDailyVisitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDailyVisitLogic {
	return &GetUserDailyVisitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户日浏览量
func (l *GetUserDailyVisitLogic) GetUserDailyVisit(in *websiterpc.EmptyReq) (*websiterpc.UserDailyVisitRsp, error) {
	result, err := l.svcCtx.TVisitHistoryModel.FindList(l.ctx, 1, 100, "id desc", "")
	if err != nil {
		return nil, err
	}

	var list []*websiterpc.UserVisit
	for _, v := range result {
		uv := &websiterpc.UserVisit{
			Date:      v.CreatedAt.Format(time.DateTime),
			ViewCount: v.ViewsCount,
		}
		list = append(list, uv)
	}

	return &websiterpc.UserDailyVisitRsp{
		List: list,
	}, nil
}
