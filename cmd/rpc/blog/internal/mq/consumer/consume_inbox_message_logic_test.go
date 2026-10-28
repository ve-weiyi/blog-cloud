package consumer

import (
	"context"
	"encoding/json"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/model"

	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/messaging/notifyx"
)

// insertWindow 放大"插入耗时"，让"事件早于落库"这类回归稳定可观测
const insertWindow = 100 * time.Millisecond

// 只覆写用得到的方法；其余靠嵌入接口兜底（一旦被用到即 nil panic，测试会立刻暴露）
type fakeRecordModel struct {
	model.TNotifyRecordModel
	existing []*model.TNotifyRecord
	inserted atomic.Bool
}

func (f *fakeRecordModel) WithTx(*gorm.DB) model.TNotifyRecordModel { return f }

func (f *fakeRecordModel) FindAll(_ context.Context, _ string, _ ...interface{}) ([]*model.TNotifyRecord, error) {
	return f.existing, nil
}

func (f *fakeRecordModel) InsertBatch(_ context.Context, in ...*model.TNotifyRecord) (int64, error) {
	time.Sleep(insertWindow)
	f.inserted.Store(true)
	return int64(len(in)), nil
}

// arrival 事件到达时的现场：一并记录"那一刻投递记录是否已落库"，
// 否则只能在 Consume 返回后检查，而那时插入必然已完成，断言就失去意义
type arrival struct {
	event     notifyx.Event
	persisted bool
}

// 订阅通知事件 channel，返回事件到达现场的通道
func subscribeNotifyEvents(t *testing.T, rdb *redis.Client, rec *fakeRecordModel) <-chan arrival {
	t.Helper()
	sub := rdb.Subscribe(context.Background(), cachekey.NotifyEventChannel)
	t.Cleanup(func() { _ = sub.Close() })

	received := make(chan arrival, 1)
	go func() {
		for msg := range sub.Channel() {
			if event, err := notifyx.Decode(msg.Payload); err == nil {
				received <- arrival{event: event, persisted: rec.inserted.Load()}
				return
			}
		}
	}()
	return received
}

// newTestLogic 用 sqlmock 驱动真实 GORM，覆盖事务边界与加锁读取。
func newTestLogic(t *testing.T, msg *model.TNotifyMessage, rec *fakeRecordModel) (*ConsumeInboxMessageLogic, <-chan arrival) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}

	var targetIds any
	if msg.TargetIds != nil {
		targetIds = *msg.TargetIds
	}

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .* FROM `t_notify_message` .*FOR UPDATE").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "title", "category", "level", "target_type", "target_ids", "status", "content",
		}).AddRow(msg.Id, msg.Title, msg.Category, msg.Level, msg.TargetType, targetIds, msg.Status, msg.Content))
	mock.ExpectCommit()

	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未满足的 SQL 期望: %v", err)
		}
	})

	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })

	svcCtx := &svc.ServiceContext{
		GormDB:             gdb,
		Redis:              rdb,
		TNotifyRecordModel: rec,
	}
	return NewConsumeInboxMessageLogic(svcCtx), subscribeNotifyEvents(t, rdb, rec)
}

func userIDsMsg(ids string) *model.TNotifyMessage {
	raw, _ := json.Marshal(strings.Split(ids, ","))
	targetIds := string(raw)
	return &model.TNotifyMessage{
		Id:         7,
		Title:      "维护通知",
		Category:   "maintenance",
		Level:      "info",
		TargetType: "user_ids",
		TargetIds:  &targetIds,
		Status:     "published",
	}
}

// 关键不变量：事件必须在投递记录落库**之后**发出。
// 反过来的话，客户端收到信号回查未读会查不到——那条通知要等下一次信号才可见。
func TestConsumePublishesOnlyAfterRecordsPersisted(t *testing.T) {
	rec := &fakeRecordModel{}
	logic, received := newTestLogic(t, userIDsMsg("u1,u2"), rec)

	if err := logic.Consume(context.Background(), &mq.InboxMessageEvent{MessageId: 7}); err != nil {
		t.Fatalf("Consume: %v", err)
	}

	select {
	case got := <-received:
		if got.event.MessageId != 7 || got.event.Type != notifyx.EventNotice {
			t.Fatalf("事件内容不符: %+v", got.event)
		}
		if !got.persisted {
			t.Fatal("不变量被破坏：事件到达时投递记录尚未落库（客户端回查会查不到）")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("未收到通知事件")
	}
}

// 事件从发出到被消费之间，消息可能已被撤回。
// 撤回路径只回收"当时已存在"的未读记录，所以消费者若照常展开，会为一条已撤回的消息
// 新插入一批未读记录——它们不会被任何后续动作回收，公告照常出现在所有人的未读列表里。
func TestConsumeSkipsRevokedMessage(t *testing.T) {
	msg := userIDsMsg("u1,u2")
	msg.Status = "revoked"

	rec := &fakeRecordModel{}
	logic, received := newTestLogic(t, msg, rec)

	if err := logic.Consume(context.Background(), &mq.InboxMessageEvent{MessageId: 7}); err != nil {
		t.Fatalf("Consume: %v", err)
	}

	if rec.inserted.Load() {
		t.Fatal("已撤回的消息不应产生投递记录")
	}
	select {
	case got := <-received:
		t.Fatalf("已撤回的消息不应发事件: %+v", got.event)
	case <-time.After(300 * time.Millisecond):
	}
}

// 没有新收件人（全部已投递）时不发事件，避免无谓的客户端回查
func TestConsumeSkipsEventWhenNothingNew(t *testing.T) {
	rec := &fakeRecordModel{existing: []*model.TNotifyRecord{
		{Recipient: "u1"},
		{Recipient: "u2"},
	}}
	logic, received := newTestLogic(t, userIDsMsg("u1,u2"), rec)

	if err := logic.Consume(context.Background(), &mq.InboxMessageEvent{MessageId: 7}); err != nil {
		t.Fatalf("Consume: %v", err)
	}

	select {
	case got := <-received:
		t.Fatalf("无新投递时不应发事件: %+v", got.event)
	case <-time.After(300 * time.Millisecond):
	}
}
