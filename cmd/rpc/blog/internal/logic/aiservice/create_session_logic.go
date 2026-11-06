package aiservicelogic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/airpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/model"
)

type CreateSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建会话
func (l *CreateSessionLogic) CreateSession(in *airpc.CreateSessionRequest) (*airpc.CreateSessionResponse, error) {
	uid, err := requireUserID(l.ctx)
	if err != nil {
		return nil, err
	}

	title := truncateRunes(strings.TrimSpace(in.Title), maxTitleRunes)
	if title == "" {
		title = "新对话"
	}

	session := &model.TAiSession{
		UserId: uid,
		Title:  title,
		Model:  strings.TrimSpace(in.Model),
	}
	if _, err := l.svcCtx.TAiSessionModel.Insert(l.ctx, session); err != nil {
		return nil, err
	}

	return &airpc.CreateSessionResponse{Session: convertSession(session)}, nil
}
