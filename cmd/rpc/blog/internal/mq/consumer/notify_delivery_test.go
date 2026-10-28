package consumer

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/model"
	mqx "github.com/ve-weiyi/vkit/adapter/mqx"
)

type fakeDeliveryModel struct {
	model.TNotifyRecordModel
	existing  *model.TNotifyRecord
	findErr   error
	insertErr error
	findCalls int
}

func (f *fakeDeliveryModel) FindOne(_ context.Context, _ string, _ string, _ ...interface{}) (*model.TNotifyRecord, error) {
	f.findCalls++
	if f.findErr != nil {
		return nil, f.findErr
	}
	if f.existing == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return f.existing, nil
}

func (f *fakeDeliveryModel) Insert(_ context.Context, _ *model.TNotifyRecord) (int64, error) {
	return 0, f.insertErr
}

type fakeChannel struct {
	calls    []string
	sendErr  error
	canRetry bool
}

func (c *fakeChannel) name() string      { return "email" }
func (c *fakeChannel) recipient() string { return "a@b.c" }

func (c *fakeChannel) resolve(context.Context) (string, string, error) {
	c.calls = append(c.calls, "resolve")
	return "tpl", "body", nil
}

func (c *fakeChannel) send(context.Context, string, string) error {
	c.calls = append(c.calls, "send")
	return c.sendErr
}

func (c *fakeChannel) retryable(error) bool { return c.canRetry }

// fakeDeliveryRecords 记录写入与更新，供 delivery 编排测试。
type fakeDeliveryRecords struct {
	model.TNotifyRecordModel
	existing *model.TNotifyRecord
	inserts  int
	updates  []map[string]interface{}
}

func (f *fakeDeliveryRecords) FindOne(context.Context, string, string, ...interface{}) (*model.TNotifyRecord, error) {
	if f.existing == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return f.existing, nil
}

func (f *fakeDeliveryRecords) Insert(_ context.Context, in *model.TNotifyRecord) (int64, error) {
	f.inserts++
	in.Id = int64(f.inserts)
	return 1, nil
}

func (f *fakeDeliveryRecords) UpdateFields(_ context.Context, fields map[string]interface{}, _ string, _ ...interface{}) (int64, error) {
	f.updates = append(f.updates, fields)
	return 1, nil
}

func TestDeliveryRun(t *testing.T) {
	ctx := context.Background()

	t.Run("已发送的记录不再投递", func(t *testing.T) {
		ch := &fakeChannel{}
		m := &fakeDeliveryRecords{existing: &model.TNotifyRecord{Id: 1, Status: enums.NotifyRecordStatusSent}}

		if err := (&delivery{records: m, ch: ch, bizID: "biz-1"}).run(ctx); err != nil {
			t.Fatalf("run() error = %v", err)
		}
		if len(ch.calls) != 0 {
			t.Fatalf("calls = %v; want none", ch.calls)
		}
		if len(m.updates) != 0 {
			t.Fatalf("updates = %v; want none", m.updates)
		}
	})

	t.Run("首次投递写入记录并回写成功", func(t *testing.T) {
		ch := &fakeChannel{}
		m := &fakeDeliveryRecords{}

		if err := (&delivery{records: m, ch: ch, bizID: "biz-1"}).run(ctx); err != nil {
			t.Fatalf("run() error = %v", err)
		}
		if m.inserts != 1 {
			t.Fatalf("inserts = %d; want 1", m.inserts)
		}
		if len(ch.calls) != 2 || ch.calls[0] != "resolve" || ch.calls[1] != "send" {
			t.Fatalf("calls = %v; want [resolve send]", ch.calls)
		}
		if got := lastStatus(m.updates); got != enums.NotifyRecordStatusSent {
			t.Fatalf("status = %v; want sent", got)
		}
	})

	t.Run("不可重试的失败被确认", func(t *testing.T) {
		ch := &fakeChannel{sendErr: errors.New("invalid email"), canRetry: false}
		m := &fakeDeliveryRecords{}

		if err := (&delivery{records: m, ch: ch, bizID: "biz-1"}).run(ctx); err != nil {
			t.Fatalf("run() error = %v; want nil", err)
		}
		if got := lastStatus(m.updates); got != enums.NotifyRecordStatusFailed {
			t.Fatalf("status = %v; want failed", got)
		}
	})

	t.Run("可重试的失败向上返回", func(t *testing.T) {
		boom := errors.New("timeout")
		ch := &fakeChannel{sendErr: boom, canRetry: true}
		m := &fakeDeliveryRecords{}

		if err := (&delivery{records: m, ch: ch, bizID: "biz-1"}).run(ctx); !errors.Is(err, boom) {
			t.Fatalf("run() error = %v; want %v", err, boom)
		}
	})

	t.Run("缺少事件ID判为永久失败", func(t *testing.T) {
		ch := &fakeChannel{}

		err := (&delivery{records: &fakeDeliveryRecords{}, ch: ch, bizID: ""}).run(ctx)
		if !mqx.IsPermanent(err) {
			t.Fatalf("run() error = %v; want permanent", err)
		}
		if len(ch.calls) != 0 {
			t.Fatalf("calls = %v; want none", ch.calls)
		}
	})
}

// lastStatus 返回最后一次状态回写的内容。
func lastStatus(updates []map[string]interface{}) interface{} {
	if len(updates) == 0 {
		return nil
	}
	return updates[len(updates)-1]["status"]
}

func TestFindDelivery(t *testing.T) {
	t.Run("empty biz id skips query", func(t *testing.T) {
		m := &fakeDeliveryModel{findErr: errors.New("should not be called")}
		got, err := findDelivery(context.Background(), m, "", "email")
		if err != nil || got != nil {
			t.Fatalf("findDelivery() = %v, %v; want nil, nil", got, err)
		}
		if m.findCalls != 0 {
			t.Fatalf("findCalls = %d; want 0", m.findCalls)
		}
	})

	t.Run("not found", func(t *testing.T) {
		got, err := findDelivery(context.Background(), &fakeDeliveryModel{}, "biz-1", "email")
		if err != nil || got != nil {
			t.Fatalf("findDelivery() = %v, %v; want nil, nil", got, err)
		}
	})

	t.Run("found", func(t *testing.T) {
		want := &model.TNotifyRecord{Id: 7}
		got, err := findDelivery(context.Background(), &fakeDeliveryModel{existing: want}, "biz-1", "email")
		if err != nil || got != want {
			t.Fatalf("findDelivery() = %v, %v; want %v, nil", got, err, want)
		}
	})
}

func TestCreateDelivery(t *testing.T) {
	t.Run("insert succeeds", func(t *testing.T) {
		in := &model.TNotifyRecord{BizId: "biz-1", Channel: "email"}
		got, err := createDelivery(context.Background(), &fakeDeliveryModel{}, in)
		if err != nil || got != in {
			t.Fatalf("createDelivery() = %v, %v; want %v, nil", got, err, in)
		}
	})

	t.Run("duplicate returns existing", func(t *testing.T) {
		existing := &model.TNotifyRecord{Id: 9, BizId: "biz-1", Channel: "email"}
		m := &fakeDeliveryModel{existing: existing, insertErr: gorm.ErrDuplicatedKey}
		got, err := createDelivery(context.Background(), m, &model.TNotifyRecord{BizId: "biz-1", Channel: "email"})
		if err != nil || got != existing {
			t.Fatalf("createDelivery() = %v, %v; want %v, nil", got, err, existing)
		}
	})

	t.Run("duplicate but not found returns error", func(t *testing.T) {
		m := &fakeDeliveryModel{insertErr: gorm.ErrDuplicatedKey}
		_, err := createDelivery(context.Background(), m, &model.TNotifyRecord{BizId: "biz-1", Channel: "email"})
		if !errors.Is(err, gorm.ErrDuplicatedKey) {
			t.Fatalf("createDelivery() error = %v; want duplicated key", err)
		}
	})

	t.Run("other insert error", func(t *testing.T) {
		boom := errors.New("boom")
		m := &fakeDeliveryModel{insertErr: boom}
		_, err := createDelivery(context.Background(), m, &model.TNotifyRecord{BizId: "biz-1", Channel: "email"})
		if !errors.Is(err, boom) {
			t.Fatalf("createDelivery() error = %v; want boom", err)
		}
	})
}
