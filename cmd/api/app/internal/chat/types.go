package chat

const (
	MessageTypeOnline      = "online"
	MessageTypeGreeting    = "greeting"
	MessageTypeHistory     = "history"
	MessageTypeSend        = "send"
	MessageTypeEdit        = "edit"
	MessageTypeMessage     = "message"
	MessageTypeHistoryMore = "history_more"
)

type MessageEvent struct {
	Type      string `json:"type"` // 事件类型 online/greeting/history/send/edit/message
	Data      string `json:"data"`
	TimeStamp int64  `json:"timestamp"`
}

type GreetingMessageEvent struct {
	Content  string `json:"content"` // 欢迎消息内容
	IpSource string `json:"ip_source"`
}

type OnlineMessageEvent struct {
	Online bool   `json:"online"` // 是否在线
	Count  int64  `json:"count"`  // 在线人数
	Tips   string `json:"tips"`   // 提示消息
}

type HistoryMessageEvent struct {
	List  []*ChatMessageEvent `json:"list"`  // 消息列表
	Page  int64               `json:"page"`  // 当前页码
	Size  int64               `json:"size"`  // 每页条数
	Total int64               `json:"total"` // 总条数
	// HasMore 由纯查询的行数判定，不按过滤后的条数推断：已删除记录也占一页名额
	HasMore bool `json:"has_more"`
	// NextBefore 是下一页的游标，取本页最早的原始行（含已删除）的 created_at
	NextBefore int64 `json:"next_before"`
}

// HistoryMoreMessageEvent 是客户端拉取更早消息的请求。
type HistoryMoreMessageEvent struct {
	Before int64 `json:"before"` // 以最早一条的 created_at 为截止时间
	Limit  int   `json:"limit"`  // 1-50
}

type ChatMessageEvent struct {
	Id          int64  `json:"id"`
	UserId      string `json:"user_id"`
	DeviceId    string `json:"device_id"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	IpSource    string `json:"ip_source"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Status      int64  `json:"status"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	ClientMsgId string `json:"client_msg_id,omitempty"` // 仅发送回显使用，不落库
}

type SendMessageEvent struct {
	Type        string `json:"type"` // 消息类型；当前只支持 text
	Content     string `json:"content"`
	ClientMsgId string `json:"client_msg_id,omitempty"`
}

type EditMessageEvent struct {
	Id        int64  `json:"id"`
	Type      string `json:"type"` // 消息类型 text/image/file/audio/video
	Content   string `json:"content"`
	Status    int64  `json:"status"` // 消息状态 0-正常 1-已编辑 2-已撤回 3-已删除
	UpdatedAt int64  `json:"updated_at"`
}
