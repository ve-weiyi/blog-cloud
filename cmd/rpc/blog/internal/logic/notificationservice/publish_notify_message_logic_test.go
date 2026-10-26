package notificationservicelogic

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/outbox"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/model"
)

func newNotifyTestLogic(t *testing.T) (*PublishNotifyMessageLogic, sqlmock.Sqlmock) {
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

	svcCtx := &svc.ServiceContext{
		GormDB:              db,
		TNotifyMessageModel: model.NewTNotifyMessageModel(db),
		Outbox:              outbox.NewGORMStore(db),
	}
	return NewPublishNotifyMessageLogic(context.Background(), svcCtx), mock
}

func TestPublishNotifyMessageWritesOutboxInTransaction(t *testing.T) {
	logic, mock := newNotifyTestLogic(t)

	mock.ExpectQuery("SELECT .*t_notify_message.*").
		WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(int64(1), "draft"))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE .*t_notify_message.*").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO .*t_outbox_event.*").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	resp, err := logic.PublishNotifyMessage(&notificationrpc.PublishNotifyMessageRequest{Id: 1})
	if err != nil {
		t.Fatalf("PublishNotifyMessage: %v", err)
	}
	if resp == nil || !resp.Success {
		t.Fatalf("resp = %+v, want success", resp)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestPublishNotifyMessageRollsBackWhenAlreadyPublished(t *testing.T) {
	logic, mock := newNotifyTestLogic(t)

	mock.ExpectQuery("SELECT .*t_notify_message.*").
		WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(int64(1), "draft"))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE .*t_notify_message.*").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	resp, err := logic.PublishNotifyMessage(&notificationrpc.PublishNotifyMessageRequest{Id: 1})
	if err != nil {
		t.Fatalf("PublishNotifyMessage: %v", err)
	}
	if resp == nil || resp.Success {
		t.Fatalf("resp = %+v, want success=false", resp)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
