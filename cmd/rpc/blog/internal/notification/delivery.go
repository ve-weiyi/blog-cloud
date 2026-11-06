package notification

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/model"
)

// writeDeliveries 为需要外部渠道的条目写待发送记录。
//
// 站在「本批收件人对应的全部条目」上写，而不是「本次新插入的条目」：重放时插入会被
// 去重键吞掉，只认新行的话，上一批漏写的那条投递就永远补不上。
func (s *Service) writeDeliveries(ctx context.Context, tx *gorm.DB, plans []Entry, address map[string]recipient, ids map[string]int64, now time.Time) error {
	recs := make([]*model.TNotificationDelivery, 0, len(plans))
	for _, p := range plans {
		nid := ids[key(p.RecipientID, p.DedupeKey)]
		if nid == 0 {
			continue
		}
		for _, ch := range p.External {
			addr := deliveryAddress(ch, address[p.RecipientID])
			if addr == "" {
				continue
			}
			recs = append(recs, &model.TNotificationDelivery{
				NotificationId: nid,
				Channel:        string(ch),
				Recipient:      addr,
				Status:         DeliveryPending,
				CreatedAt:      now,
			})
		}
	}
	if len(recs) == 0 {
		return nil
	}
	if _, err := s.deliveries.WithTx(tx).InsertIgnoreBatch(ctx, recs...); err != nil {
		return fmt.Errorf("notification: insert deliveries: %w", err)
	}
	return nil
}

// deliveryAddress 取某个渠道的投递地址；站内信由条目本身承载，不产生投递。
func deliveryAddress(ch notify.Channel, r recipient) string {
	switch ch {
	case notify.ChannelEmail:
		return r.Email
	case notify.ChannelSms:
		return r.Mobile
	default:
		return ""
	}
}
