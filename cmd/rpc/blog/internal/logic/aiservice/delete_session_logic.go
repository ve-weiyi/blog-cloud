package aiservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/airpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
)

type DeleteSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionLogic {
	return &DeleteSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除会话
func (l *DeleteSessionLogic) DeleteSession(in *airpc.DeleteSessionRequest) (*airpc.DeleteSessionResponse, error) {
	session, err := findOwnedSession(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}

	var rows int64
	err = l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		if _, err := l.svcCtx.TAiMessageModel.WithTx(tx).DeleteBatch(l.ctx, "session_id = ?", session.Id); err != nil {
			return err
		}
		rows, err = l.svcCtx.TAiSessionModel.WithTx(tx).Delete(l.ctx, session.Id)
		return err
	})
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, "会话不存在")
	}

	return &airpc.DeleteSessionResponse{SuccessCount: rows}, nil
}
