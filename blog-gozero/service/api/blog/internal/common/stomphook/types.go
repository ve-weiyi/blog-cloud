package stomphook

const (
	MessageTypeOnline   = "online"
	MessageTypeGreeting = "greeting"
	MessageTypeHistory  = "history"
	MessageTypeSend     = "send"
	MessageTypeEdit     = "edit"
	MessageTypeMessage  = "message"
)

type MessageEvent struct {
	Type      string `json:"type"` // 事件类型 online/greeting/history/send/edit/message
	Data      string `json:"data"`
	TimeStamp int64  `json:"timestamp"`
}

type GreetingMessageEvent struct {
	Content   string `json:"content"` // 欢迎消息内容
	IpAddress string `json:"ip_address"`
	IpSource  string `json:"ip_source"`
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
}

type ChatMessageEvent struct {
	Id         int64  `json:"id"`
	UserId     string `json:"user_id"`
	TerminalId string `json:"terminal_id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	IpAddress  string `json:"ip_address"`
	IpSource   string `json:"ip_source"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	Status     int64  `json:"status"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

type SendMessageEvent struct {
	Type    string `json:"type"` // 消息类型 text/image/file/audio/video
	Content string `json:"content"`
}

type EditMessageEvent struct {
	Id        int64  `json:"id"`
	Type      string `json:"type"` // 消息类型 text/image/file/audio/video
	Content   string `json:"content"`
	Status    int64  `json:"status"` // 消息状态 0-正常 1-已编辑 2-已撤回 3-已删除
	UpdatedAt int64  `json:"updated_at"`
}
