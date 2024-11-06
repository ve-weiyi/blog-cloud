package constant

// ChatTypeEnum 代表聊天消息的类型。
type ChatTypeEnum int

// 定义ChatTypeEnum的可能值。
const (
	OnlineCount   = iota + 1 // 在线人数
	HistoryRecord            // 历史记录
	SendMessage              // 发送消息
	RecallMessage            // 撤回消息
	VoiceMessage             // 语音消息
	HeartBeat                // 心跳消息
)
