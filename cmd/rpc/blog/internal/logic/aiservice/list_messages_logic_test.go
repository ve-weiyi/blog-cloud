package aiservicelogic

import (
	"context"
	"strings"
	"testing"

	"github.com/ve-weiyi/blog-cloud/model"
)

type fakeMessageLister struct {
	records    []*model.TAiMessage
	sorts      string
	conditions string
	args       []interface{}
}

func (f *fakeMessageLister) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) ([]*model.TAiMessage, int64, error) {
	f.sorts = sorts
	f.conditions = conditions
	f.args = args
	return f.records, int64(len(f.records)), nil
}

// after == 0 时按 id 倒序查询最近 N 条，再翻回正序。
func TestLoadSessionMessagesLatest(t *testing.T) {
	lister := &fakeMessageLister{records: []*model.TAiMessage{{Id: 20}, {Id: 19}, {Id: 18}}}

	records, err := loadSessionMessages(context.Background(), lister, 1, 0, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lister.sorts != "id desc" {
		t.Fatalf("sorts = %q, want id desc", lister.sorts)
	}
	if len(records) != 3 || records[0].Id != 18 || records[2].Id != 20 {
		t.Fatalf("records should be reversed to ascending, got %+v", ids(records))
	}
}

// after > 0 时按时间条件升序取增量。
func TestLoadSessionMessagesAfter(t *testing.T) {
	lister := &fakeMessageLister{records: []*model.TAiMessage{{Id: 1}, {Id: 2}}}

	records, err := loadSessionMessages(context.Background(), lister, 1, 1700000000000, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lister.sorts != "id asc" {
		t.Fatalf("sorts = %q, want id asc", lister.sorts)
	}
	if !strings.Contains(lister.conditions, "created_at >") {
		t.Fatalf("conditions = %q", lister.conditions)
	}
	if len(records) != 2 || records[0].Id != 1 {
		t.Fatalf("unexpected records: %+v", ids(records))
	}
}

func ids(records []*model.TAiMessage) []int64 {
	out := make([]int64, 0, len(records))
	for _, record := range records {
		out = append(out, record.Id)
	}
	return out
}
