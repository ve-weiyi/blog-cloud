package chatservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/chatrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/chat"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type UpdateChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateChatLogic {
	return &UpdateChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新聊天记录
func (l *UpdateChatLogic) UpdateChat(in *chatrpc.UpdateChatRequest) (*chatrpc.UpdateChatResponse, error) {
	uid, _ := metax.GetUserIdFromCtx(l.ctx)
	deviceId, _ := metax.GetDeviceIdFromCtx(l.ctx)
	if uid == "" && deviceId == "" {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "缺少用户身份")
	}

	content := in.Content
	messageType := in.Type
	switch in.Status {
	case enums.ChatStatusRecall:
		// 撤回文案由服务端固定，客户端不能把它替换成任意内容
		content = "[消息已撤回]"
		messageType = chat.TypeText
	case enums.ChatStatusNormal, enums.ChatStatusEdited:
		if !chat.ValidType(messageType) {
			return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "不支持的消息类型")
		}
		normalized, err := chat.NormalizeContent(content)
		if err != nil {
			return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "消息内容不合法")
		}
		content = normalized
	default:
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "聊天状态不合法")
	}

	conditions := "id = ?"
	args := []any{in.Id}
	if uid != "" {
		conditions += " AND user_id = ?"
		args = append(args, uid)
	} else {
		conditions += " AND device_id = ?"
		args = append(args, deviceId)
	}

	rows, err := l.svcCtx.TChatModel.UpdateFields(l.ctx, map[string]interface{}{
		"type":    messageType,
		"content": content,
		"status":  in.Status,
	}, conditions, args...)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, bizerr.NewBizError(bizcode.CodeNoPermission, "无权修改该聊天记录")
	}

	return &chatrpc.UpdateChatResponse{Success: true}, nil
}
