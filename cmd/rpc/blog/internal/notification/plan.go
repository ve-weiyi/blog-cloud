package notification

import (
	"fmt"
	"strings"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

// 受众类型，与 t_notification_event.audience_type 的取值一致。
const (
	AudienceAll     = "all"
	AudienceUserIDs = "user_ids"
)

// ValidAudience 判断受众类型是否合法。空值视为非法：漏填时直接报错，
// 而不是默认成全站广播。
func ValidAudience(audienceType string) bool {
	return audienceType == AudienceAll || audienceType == AudienceUserIDs
}

// NormalizeRecipients 清洗收件人：去空白、去空值、去重，并剔除触发者本人。
//
// 顺序保持首次出现的次序，便于分页游标在同一批内稳定推进。
func NormalizeRecipients(ids []string, actorID string) []string {
	if len(ids) == 0 {
		return nil
	}

	out := make([]string, 0, len(ids))
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || id == actorID {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

// GroupKey 计算展示分组键：按对象分组的场景合并同一条目下的互动，其余场景一条一组。
//
// 不分组场景用事件标识而不是条目 id，是因为条目 id 要插入后才存在，
// 而组键必须在插入前算好，才能连去重键一起写进同一行。
func GroupKey(scene notify.Scene, objectType string, objectID, eventID int64) string {
	if scene.Grouping == notify.GroupingObject {
		return fmt.Sprintf("%s:%s:%d", scene.Key, objectType, objectID)
	}
	return fmt.Sprintf("%s:e:%d", scene.Key, eventID)
}

// DedupeKey 计算幂等键：可分组场景含触发者，同一个人对同一对象只通知一次；
// 不可分组场景与组键相同，重复扇出由它吞掉。
func DedupeKey(scene notify.Scene, groupKey, actorID string) string {
	if scene.Grouping == notify.GroupingObject {
		return groupKey + ":" + actorID
	}
	return groupKey
}

// InteractionKey 定位「某触发者对某对象的一次互动」写下的条目去重键。
//
// 撤回这类条目时只有对象与触发者可用，事件 id 未知；可分组场景的组键不含事件 id，
// 所以这里的 0 只是占位，不参与结果。
func InteractionKey(scene notify.Scene, objectType string, objectID int64, actorID string) string {
	return DedupeKey(scene, GroupKey(scene, objectType, objectID, 0), actorID)
}

// Channels 返回某个场景对某个收件人实际应走的渠道。
//
// prefs 是该收件人在这个偏好分组下的显式开关，键为渠道；缺项表示沿用场景默认。
//
// 两条规则：
//   - 强制分组不看偏好：公告与账号安全不可退订；
//   - 站内信是锚点渠道：关闭它等于关闭该分组的全部渠道，因为条目同时是外部投递的载体。
func Channels(scene notify.Scene, prefs map[notify.Channel]bool) []notify.Channel {
	if scene.Mandatory {
		return scene.Channels
	}

	if enabled, ok := prefs[notify.ChannelInbox]; ok && !enabled {
		return nil
	}

	out := make([]notify.Channel, 0, len(scene.Channels))
	for _, ch := range scene.Channels {
		if ch != notify.ChannelInbox {
			if enabled, ok := prefs[ch]; ok && !enabled {
				continue
			}
		}
		out = append(out, ch)
	}
	return out
}

// Entry 是一条待写入的收件箱条目。
type Entry struct {
	RecipientID string
	GroupKey    string
	DedupeKey   string
	// External 只含站内信以外的渠道；站内信由条目本身承载，不进投递表。
	External []notify.Channel
}

// PlanEntries 把一批收件人换算成待写入的条目。
//
// 偏好判定在这里完成：渠道全关的收件人直接跳过，不产生条目 —— 关掉站内信等于关掉该分组全渠道，
// 因为条目同时是外部投递的载体。站内信本身不落投递表，只保留外部渠道。
func PlanEntries(scene notify.Scene, objectType string, objectID, eventID int64, actorID string, recipients []string, prefs map[string]map[notify.Channel]bool) []Entry {
	groupKey := GroupKey(scene, objectType, objectID, eventID)

	out := make([]Entry, 0, len(recipients))
	for _, uid := range recipients {
		channels := Channels(scene, prefs[uid])
		if len(channels) == 0 {
			continue
		}
		out = append(out, Entry{
			RecipientID: uid,
			GroupKey:    groupKey,
			DedupeKey:   DedupeKey(scene, groupKey, actorID),
			External:    externalChannels(channels),
		})
	}
	return out
}

// externalChannels 去掉站内信：它由条目本身承载，不产生投递记录。
func externalChannels(channels []notify.Channel) []notify.Channel {
	out := make([]notify.Channel, 0, len(channels))
	for _, ch := range channels {
		if ch != notify.ChannelInbox {
			out = append(out, ch)
		}
	}
	return out
}
