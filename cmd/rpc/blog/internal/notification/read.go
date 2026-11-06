package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/model"
)

// maxPageSize 是收件箱列表单页条目的上限。
const maxPageSize = 100

// InboxGroup 是收件箱里的一个展示单元。
type InboxGroup struct {
	GroupKey        string
	Scene           string
	Title           string
	Content         string
	Link            string
	Level           string
	ActorCount      int64
	LatestActorName string
	LatestAt        int64
	Unread          bool
}

// ListGroups 按分组分页返回收件箱。
//
// groups 为空表示全部偏好分组；服务内部再展开成场景集合，调用方不必自己转换。
// unreadOnly 为真时，组大小与最新触发者都只统计未读条目。
func (s *Service) ListGroups(ctx context.Context, recipientID string, groups []notify.Group, unreadOnly bool, page, size int, lang string) ([]InboxGroup, int64, error) {
	scenes := notify.ScenesOfGroups(groups)
	if len(scenes) == 0 {
		return nil, 0, nil
	}
	if size <= 0 || size > maxPageSize {
		size = maxPageSize
	}
	if page <= 0 {
		page = 1
	}

	// 每次都从干净的条件重建，避免在同一条查询链上叠加 Select
	filter := func() *gorm.DB {
		q := s.db.WithContext(ctx).Table(s.notifs.TableName()).
			Where("recipient_id = ? AND scene IN ?", recipientID, scenes)
		if unreadOnly {
			q = q.Where("read_at IS NULL")
		}
		return q
	}

	var total int64
	if err := filter().Select("COUNT(DISTINCT group_key)").Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("notification: count groups: %w", err)
	}

	var rows []struct {
		GroupKey  string
		LatestId  int64
		LatestAt  time.Time
		Cnt       int64
		UnreadCnt int64
	}
	// 按 MAX(id) 分页：id 单调，与「最新一条」等价，且天然没有并列排序键
	err := filter().
		Select("group_key AS group_key, MAX(id) AS latest_id, MAX(created_at) AS latest_at, COUNT(*) AS cnt, SUM(read_at IS NULL) AS unread_cnt").
		Group("group_key").
		Order("latest_id DESC").
		Offset((page - 1) * size).
		Limit(size).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("notification: list groups: %w", err)
	}
	if len(rows) == 0 {
		return nil, total, nil
	}

	latestIDs := make([]int64, 0, len(rows))
	for _, r := range rows {
		latestIDs = append(latestIDs, r.LatestId)
	}
	latest, err := s.loadEntries(ctx, latestIDs)
	if err != nil {
		return nil, 0, err
	}

	// 整页一次性取事件与公告，避免每组各查一次
	events, err := s.loadEvents(ctx, latest)
	if err != nil {
		return nil, 0, err
	}
	announcements, err := s.loadAnnouncements(ctx, latest, events)
	if err != nil {
		return nil, 0, err
	}

	groupsOut := make([]InboxGroup, 0, len(rows))
	for _, r := range rows {
		entry := latest[r.LatestId]
		if entry == nil {
			continue
		}
		ev := events[entry.EventId]
		if ev == nil {
			continue
		}
		groupsOut = append(groupsOut, project(entry, ev, announcements, r.Cnt, r.UnreadCnt, r.LatestAt, lang))
	}
	return groupsOut, total, nil
}

// CountUnread 返回未读分组总数与按偏好分组的分别计数。
//
// 计数按展示组（收件箱里的一组）算，但对外以偏好分组为键下发 —— 前端拿的是
// 「点赞 3、评论 1」这类角标，不是按场景拆开的细目。
func (s *Service) CountUnread(ctx context.Context, recipientID string, groups []notify.Group) (int64, map[notify.Group]int64, error) {
	out := make(map[notify.Group]int64, len(groups))
	scenes := notify.ScenesOfGroups(groups)
	if len(scenes) == 0 {
		return 0, out, nil
	}

	var rows []struct {
		Scene string
		Cnt   int64
	}
	err := s.db.WithContext(ctx).Table(s.notifs.TableName()).
		Select("scene AS scene, COUNT(DISTINCT group_key) AS cnt").
		Where("recipient_id = ? AND scene IN ? AND read_at IS NULL", recipientID, scenes).
		Group("scene").
		Scan(&rows).Error
	if err != nil {
		return 0, nil, fmt.Errorf("notification: count unread: %w", err)
	}

	var total int64
	for _, r := range rows {
		scene, ok := notify.Lookup(r.Scene)
		if !ok {
			continue
		}
		out[scene.Group] += r.Cnt
		total += r.Cnt
	}
	return total, out, nil
}

// loadEntries 按 id 取条目。
func (s *Service) loadEntries(ctx context.Context, ids []int64) (map[int64]*model.TNotification, error) {
	rows, err := s.notifs.FindAll(ctx, "id IN ?", ids)
	if err != nil {
		return nil, fmt.Errorf("notification: load entries: %w", err)
	}

	out := make(map[int64]*model.TNotification, len(rows))
	for _, row := range rows {
		out[row.Id] = row
	}
	return out, nil
}

// loadEvents 取这批条目对应的事件。
func (s *Service) loadEvents(ctx context.Context, entries map[int64]*model.TNotification) (map[int64]*model.TNotificationEvent, error) {
	ids := make([]int64, 0, len(entries))
	seen := make(map[int64]struct{}, len(entries))
	for _, e := range entries {
		if _, ok := seen[e.EventId]; ok {
			continue
		}
		seen[e.EventId] = struct{}{}
		ids = append(ids, e.EventId)
	}

	rows, err := s.events.FindAll(ctx, "id IN ?", ids)
	if err != nil {
		return nil, fmt.Errorf("notification: load events: %w", err)
	}

	out := make(map[int64]*model.TNotificationEvent, len(rows))
	for _, row := range rows {
		out[row.Id] = row
	}
	return out, nil
}

// loadAnnouncements 取本页公告条目的正文。
func (s *Service) loadAnnouncements(ctx context.Context, entries map[int64]*model.TNotification, events map[int64]*model.TNotificationEvent) (map[int64]*model.TAnnouncement, error) {
	ids := make([]int64, 0, len(entries))
	seen := make(map[int64]struct{}, len(entries))
	for _, e := range entries {
		if e.Scene != notify.SceneAnnouncement {
			continue
		}
		ev := events[e.EventId]
		if ev == nil {
			continue
		}
		id := decodePayload(ev.Payload).AnnouncementID
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, nil
	}

	rows, err := s.announcements.FindAll(ctx, "id IN ?", ids)
	if err != nil {
		return nil, fmt.Errorf("notification: load announcements: %w", err)
	}

	out := make(map[int64]*model.TAnnouncement, len(rows))
	for _, row := range rows {
		out[row.Id] = row
	}
	return out, nil
}

// project 渲染一个分组：公告取内容表，其余取事件载荷。
func project(entry *model.TNotification, ev *model.TNotificationEvent, announcements map[int64]*model.TAnnouncement, count, unreadCount int64, latestAt time.Time, lang string) InboxGroup {
	payload := decodePayload(ev.Payload)
	title, body, link, level := renderContent(entry, ev, announcements, count, lang)

	// 公告的载荷只存内容引用，展示「谁发的」要取发布人
	actor := payload.ActorName
	if entry.Scene == notify.SceneAnnouncement {
		if row := announcements[payload.AnnouncementID]; row != nil {
			actor = row.PublishedBy
		}
	}

	return InboxGroup{
		GroupKey:        entry.GroupKey,
		Scene:           entry.Scene,
		Title:           title,
		Content:         body,
		Link:            link,
		Level:           level,
		ActorCount:      count,
		LatestActorName: actor,
		LatestAt:        latestAt.Unix(),
		Unread:          unreadCount > 0,
	}
}

// renderContent 渲染展示文本：公告读内容表，其余按场景与载荷渲染。
func renderContent(entry *model.TNotification, ev *model.TNotificationEvent, announcements map[int64]*model.TAnnouncement, count int64, lang string) (title, body, link, level string) {
	payload := decodePayload(ev.Payload)

	if entry.Scene == notify.SceneAnnouncement {
		if row := announcements[payload.AnnouncementID]; row != nil {
			return row.Title, deref(row.Content), "", row.Level
		}
		return "公告已下线", "", "", string(notify.LevelInfo)
	}

	rendered, _ := notify.Render(entry.Scene, payload, count, lang)
	if sc, ok := notify.Lookup(entry.Scene); ok {
		level = string(sc.Level)
	}
	return rendered.Title, rendered.Body, rendered.Link, level
}

// Item 是条目级视图：管理端排障用，含单条已读时间与原始标识。
type Item struct {
	ID          int64
	EventID     int64
	RecipientID string
	Scene       string
	ActorID     string
	ActorName   string
	ObjectType  string
	ObjectID    int64
	Title       string
	Content     string
	Link        string
	ReadAt      *time.Time
	CreatedAt   time.Time
}

// ListItems 按条目分页返回某个收件人的收件箱，不做分组。
//
// 管理端排障要看单条的已读时间与原始标识；分组投影满足不了这个需求。
func (s *Service) ListItems(ctx context.Context, recipientID string, page, size int, lang string) ([]Item, int64, error) {
	if size <= 0 || size > maxPageSize {
		size = maxPageSize
	}
	if page <= 0 {
		page = 1
	}

	rows, total, err := s.notifs.FindListAndTotal(ctx, page, size, "id DESC", "recipient_id = ?", recipientID)
	if err != nil {
		return nil, 0, fmt.Errorf("notification: list items: %w", err)
	}
	if len(rows) == 0 {
		return nil, total, nil
	}

	byID := make(map[int64]*model.TNotification, len(rows))
	for _, row := range rows {
		byID[row.Id] = row
	}
	events, err := s.loadEvents(ctx, byID)
	if err != nil {
		return nil, 0, err
	}
	announcements, err := s.loadAnnouncements(ctx, byID, events)
	if err != nil {
		return nil, 0, err
	}

	out := make([]Item, 0, len(rows))
	for _, entry := range rows {
		ev := events[entry.EventId]
		if ev == nil {
			continue
		}
		title, body, link, _ := renderContent(entry, ev, announcements, 1, lang)
		out = append(out, Item{
			ID:          entry.Id,
			EventID:     entry.EventId,
			RecipientID: entry.RecipientId,
			Scene:       entry.Scene,
			ActorID:     ev.ActorId,
			ActorName:   decodePayload(ev.Payload).ActorName,
			ObjectType:  entry.ObjectType,
			ObjectID:    entry.ObjectId,
			Title:       title,
			Content:     body,
			Link:        link,
			ReadAt:      entry.ReadAt,
			CreatedAt:   entry.CreatedAt,
		})
	}
	return out, total, nil
}

func decodePayload(raw *string) notify.Payload {
	var p notify.Payload
	if raw == nil || *raw == "" {
		return p
	}
	_ = json.Unmarshal([]byte(*raw), &p)
	return p
}

func deref(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
