package messagerpclogic

import (
	"context"
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRemarkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRemarkLogic {
	return &UpdateRemarkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新留言
func (l *UpdateRemarkLogic) UpdateRemark(in *messagerpc.UpdateRemarkReq) (*messagerpc.UpdateRemarkResp, error) {
	uid, err := rpcutils.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	remark, err := l.svcCtx.TRemarkModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if remark.UserId != uid {
		return nil, fmt.Errorf("无权限操作")
	}

	remark.MessageContent = in.MessageContent
	remark.Status = in.Status
	_, err = l.svcCtx.TRemarkModel.Save(l.ctx, remark)
	if err != nil {
		return nil, err
	}

	return &messagerpc.UpdateRemarkResp{
		Remark: convertRemarkOut(remark),
	}, nil
}
