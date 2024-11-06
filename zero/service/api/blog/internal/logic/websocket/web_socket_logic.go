package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/constant"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/ipx"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/common/rediskey"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/messagerpc"
)

type WebSocketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// WebSocket消息
func NewWebSocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebSocketLogic {
	return &WebSocketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebSocketLogic) WebSocket(w http.ResponseWriter, r *http.Request) error {
	client, err := l.svcCtx.WebsocketManager.RegisterWebSocket(w, r)
	if err != nil {
		return err
	}

	onlineKey := rediskey.GetChatOnlineKey()
	// 将客户端 ID 加入 Redis Set
	l.svcCtx.Redis.ZaddCtx(l.ctx, onlineKey, time.Now().Unix(), client.Name)
	defer func() {
		l.svcCtx.Redis.ZremCtx(l.ctx, onlineKey, client.Name)
	}()

	// 有人上线，广播消息在线人数
	logx.WithContext(r.Context()).Infof("online: %v", client.Name)
	count, err := l.svcCtx.Redis.ZcardCtx(l.ctx, onlineKey)
	if err != nil {
		return err
	}

	online := l.newSocketMsg(ONLINE_COUNT, count)

	// 广播在线消息
	err = l.svcCtx.WebsocketManager.BroadcastMsg([]byte(jsonconv.AnyToJsonNE(online)))
	if err != nil {
		return err
	}

	// 获取历史记录
	history, err := l.onHistoryRecord("")
	if err != nil {
		return err
	}

	// 发送历史记录
	historyMsg := l.newSocketMsg(HISTORY_RECORD, history)

	err = client.PublishMessage(nil, []byte(jsonconv.AnyToJsonNE(historyMsg)))
	if err != nil {
		return err
	}

	// 接收消息
	receive := func(msg []byte) (err error) {
		logx.WithContext(r.Context()).Infof("receive msg: %v", string(msg))

		token := cast.ToString(r.Context().Value(constant.HeaderAuthorization))
		uid := cast.ToString(r.Context().Value(constant.HeaderUid))
		// 如果有uid,则需要校验用户是否登录
		if token != "" || uid != "" {
			_, err = l.svcCtx.TokenHolder.VerifyToken(r.Context(), token, cast.ToString(uid))
			if err != nil {
				return err
			}
		}

		var req types.WebSocketMsg
		err = json.Unmarshal(msg, &req)
		if err != nil {
			return fmt.Errorf("msg content must be :%v", jsonconv.AnyToJsonNE(req))
		}

		switch req.Cmd {
		case HEART_BEAT:
			// 心跳
			_, err := l.onHeartBeat(req.Data)
			if err != nil {
				return err
			}

			heartMsg := l.newSocketMsg(HEART_BEAT, nil)

			return client.PublishMessage(nil, []byte(jsonconv.AnyToJsonNE(heartMsg)))

		case SEND_MESSAGE:
			// 发送消息
			send, err := l.onSendMessage(req.Data)
			if err != nil {
				return err
			}

			sendMsg := l.newSocketMsg(SEND_MESSAGE, send)

			// 广播消息
			return l.svcCtx.WebsocketManager.BroadcastMsg([]byte(jsonconv.AnyToJsonNE(sendMsg)))
		case RECALL_MESSAGE:
			// 撤回消息
			rc, err := l.onRecallMessage(req.Data)
			if err != nil {
				return err
			}

			rcMsg := l.newSocketMsg(RECALL_MESSAGE, rc)
			// 广播消息
			return l.svcCtx.WebsocketManager.BroadcastMsg([]byte(jsonconv.AnyToJsonNE(rcMsg)))
		default:
			unknownMsg := l.newSocketMsg(0, "unknown message")
			// 广播消息
			return l.svcCtx.WebsocketManager.BroadcastMsg([]byte(jsonconv.AnyToJsonNE(unknownMsg)))
		}
	}

	return client.SubscribeMessage(receive)
}

// 心跳
func (l *WebSocketLogic) onHeartBeat(payload string) (data any, err error) {
	logx.WithContext(l.ctx).Infof("heartbeat: %v", payload)

	addr := cast.ToString(l.ctx.Value(constant.HeaderRemoteAddr))
	onlineKey := rediskey.GetChatOnlineKey()
	_, err = l.svcCtx.Redis.ZaddCtx(l.ctx, onlineKey, time.Now().Unix(), addr)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// 获取在线人数
func (l *WebSocketLogic) onOnlineCount(payload string) (data int, err error) {
	logx.WithContext(l.ctx).Infof("online count: %v", payload)

	onlineKey := rediskey.GetChatOnlineKey()
	count, err := l.svcCtx.Redis.ZcardCtx(l.ctx, onlineKey)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// 获取历史记录
func (l *WebSocketLogic) onHistoryRecord(payload string) (data []*types.ChatMsgResp, err error) {
	logx.WithContext(l.ctx).Infof("history record: %v", payload)

	in := &messagerpc.FindChatMessageListReq{
		After:  0,
		Before: time.Now().Unix(),
	}

	out, err := l.svcCtx.MessageRpc.FindChatMessageList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.ChatMsgResp, 0)
	for _, v := range out.List {
		msg := ConvertChatMessageTypes(v)

		list = append(list, msg)
	}

	return list, nil
}

// 发送消息
func (l *WebSocketLogic) onSendMessage(payload string) (data *types.ChatMsgResp, err error) {
	logx.WithContext(l.ctx).Infof("send message: %v", payload)

	uid := cast.ToString(l.ctx.Value(constant.HeaderUid))
	device := cast.ToString(l.ctx.Value(constant.HeaderTerminal))
	addr := cast.ToString(l.ctx.Value(constant.HeaderRemoteAddr))

	// 发送消息
	is, err := ipx.GetIpSourceByBaidu(addr)
	if err != nil {
		return nil, err
	}

	var chat types.ChatMsgReq
	err = json.Unmarshal([]byte(payload), &chat)
	if err != nil {
		return nil, err
	}

	in := &messagerpc.ChatMessageNewReq{
		Id:          0,
		UserId:      uid,
		DeviceId:    device,
		TopicId:     "",
		ReplyMsgId:  "",
		ReplyUserId: "",
		IpAddress:   addr,
		IpSource:    is,
		ChatContent: chat.ChatContent,
		Type:        chat.Type,
	}

	out, err := l.svcCtx.MessageRpc.AddChatMessage(l.ctx, in)
	if err != nil {
		return nil, err
	}

	data = ConvertChatMessageTypes(out)

	return data, nil
}

// 撤回消息
func (l *WebSocketLogic) onRecallMessage(payload string) (out *types.ChatMsgResp, err error) {
	logx.WithContext(l.ctx).Infof("recall message: %v", payload)

	return nil, nil
}

func (l *WebSocketLogic) newSocketMsg(cmd int64, content interface{}) *types.WebSocketMsg {
	device := cast.ToString(l.ctx.Value(constant.HeaderTerminal))
	addr := cast.ToString(l.ctx.Value(constant.HeaderRemoteAddr))

	msg := &types.WebSocketMsg{
		ClientId:  device,
		ClientIp:  addr,
		Timestamp: time.Now().Unix(),
		Cmd:       cmd,
		Data:      jsonconv.AnyToJsonNE(content),
	}

	return msg
}

func ConvertChatMessageTypes(in *messagerpc.ChatMessageDetails) *types.ChatMsgResp {
	return &types.ChatMsgResp{
		Id:          in.Id,
		UserId:      in.UserId,
		DeviceId:    in.DeviceId,
		Nickname:    "",
		Avatar:      "",
		ChatContent: in.ChatContent,
		IpAddress:   in.IpAddress,
		IpSource:    in.IpSource,
		Type:        in.Type,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}

const (
	ONLINE_COUNT   = 1 // 在线人数
	HISTORY_RECORD = 2 // 历史记录
	SEND_MESSAGE   = 3 // 发送消息
	RECALL_MESSAGE = 4 // 撤回消息
	HEART_BEAT     = 5 // 心跳
)

// 消息类型 1: 文本消息 2: 图片消息 3: 文件消息 4: 语音消息 5: 视频消息
const ()
