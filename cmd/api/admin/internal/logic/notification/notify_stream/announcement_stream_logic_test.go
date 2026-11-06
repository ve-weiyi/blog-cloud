package notify_stream

import (
	"testing"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/messaging/notifyx"
)

func TestToStreamEventMapsAllFields(t *testing.T) {
	in := notifyx.Event{
		Type:           notifyx.EventAnnouncementPublished,
		AnnouncementId: 5,
		Title:          "维护通知",
		Category:       "maintenance",
		Level:          "warning",
		PublishedAt:    1758000000,
	}

	got := toStreamEvent(in)

	want := types.AnnouncementStreamEvent{
		Event:          notifyx.EventAnnouncementPublished,
		AnnouncementId: 5,
		Title:          "维护通知",
		Category:       "maintenance",
		Level:          "warning",
		PublishedAt:    1758000000,
	}
	if *got != want {
		t.Fatalf("映射不符:\n got=%+v\nwant=%+v", *got, want)
	}
}
