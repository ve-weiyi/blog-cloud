package notification_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/model"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		TranslateError: true,
		// 与生产一致：模型没有 TableName 方法，靠单数命名策略与显式表名对齐
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("gorm: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return gdb, mock
}

var selectPrefs = regexp.QuoteMeta("SELECT * FROM `t_notification_preference` WHERE user_id = ?")

// 强制分组与未知分组一律不落库：前者是安全底线，后者是脏数据。
func TestUpdatePreferences_IgnoresMandatoryAndUnknown(t *testing.T) {
	gdb, mock := newMockDB(t)

	mock.ExpectQuery(selectPrefs).
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "pref_group", "channel", "enabled"}))

	svc := notification.NewService(notification.Deps{DB: gdb, Preferences: model.NewTNotificationPreferenceModel(gdb)})
	got, err := svc.UpdatePreferences(context.Background(), "u1", []notification.Preference{
		{Group: "security", Channel: "email", Enabled: false},
		{Group: "nope", Channel: "email", Enabled: false},
	})
	if err != nil {
		t.Fatalf("UpdatePreferences: %v", err)
	}

	var sawSecurity bool
	for _, p := range got {
		if p.Group == "security" {
			sawSecurity = true
			if !p.Enabled || !p.Mandatory {
				t.Errorf("强制分组不应被关闭：%+v", p)
			}
		}
	}
	if !sawSecurity {
		t.Error("返回的矩阵里应包含安全分组")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("不应发生任何写入：%v", err)
	}
}

// 普通分组的关闭会 UPSERT 落库，并回显到结果里。
func TestUpdatePreferences_WritesNonMandatory(t *testing.T) {
	gdb, mock := newMockDB(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `t_notification_preference`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(selectPrefs).
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "pref_group", "channel", "enabled"}).
			AddRow(1, "u1", "review", "email", 0))

	svc := notification.NewService(notification.Deps{DB: gdb, Preferences: model.NewTNotificationPreferenceModel(gdb)})
	got, err := svc.UpdatePreferences(context.Background(), "u1", []notification.Preference{
		{Group: "review", Channel: "email", Enabled: false},
	})
	if err != nil {
		t.Fatalf("UpdatePreferences: %v", err)
	}

	var found bool
	for _, p := range got {
		if p.Group == "review" && p.Channel == "email" {
			found = true
			if p.Enabled {
				t.Errorf("关闭邮件后应回显为 false：%+v", p)
			}
		}
	}
	if !found {
		t.Error("返回的矩阵里应包含审核分组的邮件渠道")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 未按预期执行：%v", err)
	}
}
