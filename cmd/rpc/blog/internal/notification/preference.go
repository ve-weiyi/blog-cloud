package notification

import (
	"context"
	"fmt"
	"slices"

	"gorm.io/gorm/clause"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/model"
)

// Preference 是偏好矩阵里的一格：某个分组下的某个渠道。
type Preference struct {
	Group     notify.Group
	Channel   notify.Channel
	Enabled   bool
	Mandatory bool
}

// ListPreferences 返回完整的偏好矩阵：没有显式记录的格子按场景默认回退，
// 强制分组恒为开启，客户端据此渲染开关而不是自己推断默认值。
func (s *Service) ListPreferences(ctx context.Context, userID string) ([]Preference, error) {
	rows, err := s.prefs.FindAll(ctx, "user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("notification: load preferences: %w", err)
	}

	saved := make(map[cell]bool, len(rows))
	for _, row := range rows {
		saved[cell{row.PrefGroup, row.Channel}] = row.Enabled == 1
	}

	out := make([]Preference, 0, 8)
	for _, g := range notify.Groups() {
		info, _ := notify.GroupInfoOf(g)
		for _, ch := range info.Channels {
			enabled := true
			if !info.Mandatory {
				if v, ok := saved[cell{string(g), string(ch)}]; ok {
					enabled = v
				}
			}
			out = append(out, Preference{Group: g, Channel: ch, Enabled: enabled, Mandatory: info.Mandatory})
		}
	}
	return out, nil
}

// UpdatePreferences 写入偏好并返回更新后的矩阵。
//
// 强制分组与未知分组一律忽略：mandatory 由场景目录决定，不能由客户端传参覆盖，
// 否则「安全通知不可退订」可以被绕过。
func (s *Service) UpdatePreferences(ctx context.Context, userID string, in []Preference) ([]Preference, error) {
	for _, p := range in {
		info, ok := notify.GroupInfoOf(p.Group)
		if !ok || info.Mandatory || !slices.Contains(info.Channels, p.Channel) {
			continue
		}
		if err := s.savePreference(ctx, userID, p.Group, p.Channel, p.Enabled); err != nil {
			return nil, err
		}
	}
	return s.ListPreferences(ctx, userID)
}

// savePreference 按业务唯一键 UPSERT 一格；并发下的唯一键冲突退化为更新。
func (s *Service) savePreference(ctx context.Context, userID string, group notify.Group, ch notify.Channel, enabled bool) error {
	flag := int64(0)
	if enabled {
		flag = 1
	}

	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "pref_group"}, {Name: "channel"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"enabled": flag}),
	}).Create(&model.TNotificationPreference{
		UserId:    userID,
		PrefGroup: string(group),
		Channel:   string(ch),
		Enabled:   flag,
	}).Error
}

// cell 是偏好矩阵的一格，用作已存记录的索引键。
type cell struct {
	group   string
	channel string
}
