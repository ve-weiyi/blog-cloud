package cachekey

// Notification模块缓存Key定义

const (
	// NotifyEventChannel 公告状态变化的 Pub/Sub channel。
	// 管理端网关订阅它，用于多管理员协作时刷新列表。
	NotifyEventChannel = "blog:notifyx:event"

	// NotifyUserChannelPrefix 读者端「收件箱变了」的按收件人频道前缀。
	// 信号只发给持有该用户的实例，不向全站在线连接广播。
	NotifyUserChannelPrefix = "blog:notifyx:user:"
)

// NotifyUserChannel 返回某个收件人的信号频道。
func NotifyUserChannel(userID string) string {
	return NotifyUserChannelPrefix + userID
}
