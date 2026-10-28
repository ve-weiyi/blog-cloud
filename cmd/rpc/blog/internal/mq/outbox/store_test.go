package outbox

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ve-weiyi/vkit/adapter/mqx"
)

func newTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	return db, mock
}

func TestEnqueueInTransaction(t *testing.T) {
	db, mock := newTestDB(t)
	store := NewGORMStore(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO .*t_outbox_event").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := db.Transaction(func(tx *gorm.DB) error {
		return store.Enqueue(context.Background(), tx, mqx.Message{
			ID:          "event-1",
			Destination: "test.exchange",
			RoutingKey:  "test.event",
			Body:        []byte(`{"id":1}`),
		})
	})
	if err != nil {
		t.Fatalf("Enqueue: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestClaimUsesSkipLocked(t *testing.T) {
	db, mock := newTestDB(t)
	store := NewGORMStore(db)
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*t_outbox_event.*FOR UPDATE SKIP LOCKED").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "event_id", "destination", "routing_key", "message_key",
			"headers", "payload", "status", "retry_count", "next_attempt_at",
			"locked_at", "lock_token", "last_error", "occurred_at", "created_at", "updated_at", "published_at",
		}).AddRow(
			int64(1), "event-1", "test.exchange", "test.event", "",
			`{}`, []byte(`{"id":1}`), StatusPending, 0, now,
			nil, "", "", now, now, now, nil,
		))
	mock.ExpectExec("UPDATE .*t_outbox_event").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	rows, err := store.Claim(context.Background(), 10)
	if err != nil {
		t.Fatalf("Claim: %v", err)
	}
	if len(rows) != 1 || rows[0].ID != 1 {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].LockToken == "" {
		t.Fatal("claim token is empty")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
