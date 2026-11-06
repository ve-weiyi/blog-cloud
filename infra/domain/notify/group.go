package notify

import "slices"

// GroupInfo 是偏好分组的派生元数据：渠道集合与是否强制。
type GroupInfo struct {
	// Channels 站内信恒在首位，外部渠道按名称排序，便于偏好矩阵稳定回显。
	Channels []Channel
	// Mandatory 分组内含任一强制场景即为真，整个分组不可退订。
	Mandatory bool
}

var groupInfos = buildGroupInfos()

// GroupInfoOf 返回偏好分组的元数据；未注册的分组返回 false。
func GroupInfoOf(g Group) (GroupInfo, bool) {
	info, ok := groupInfos[g]
	return info, ok
}

// buildGroupInfos 从场景目录推导分组元数据，避免调用方每次遍历 All() 重算。
func buildGroupInfos() map[Group]GroupInfo {
	infos := make(map[Group]GroupInfo, len(Groups()))
	for _, g := range Groups() {
		infos[g] = GroupInfo{Channels: []Channel{ChannelInbox}}
	}

	for _, sc := range All() {
		info := infos[sc.Group]
		info.Mandatory = info.Mandatory || sc.Mandatory
		for _, ch := range sc.Channels {
			if !slices.Contains(info.Channels, ch) {
				info.Channels = append(info.Channels, ch)
			}
		}
		infos[sc.Group] = info
	}

	for g, info := range infos {
		slices.Sort(info.Channels[1:])
		infos[g] = info
	}
	return infos
}
