package model_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/model"
)

func newTestLikeModel(t *testing.T) (model.TLikeModel, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{
		SkipDefaultTransaction: true,
		TranslateError:         true,
	})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	return model.NewTLikeModel(db), mock
}

func TestTLikeModelLikeReportsChangeFromAffectedRows(t *testing.T) {
	const insert = "INSERT INTO `t_like` .* ON DUPLICATE KEY UPDATE"

	cases := []struct {
		name        string
		affected    int64
		wantChanged bool
	}{
		{"首次点赞", 1, true},
		{"重复点赞", 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			likeModel, mock := newTestLikeModel(t)
			mock.ExpectExec(insert).WillReturnResult(sqlmock.NewResult(1, c.affected))

			changed, err := likeModel.Like(context.Background(), 1, 10, "u1")
			if err != nil {
				t.Fatalf("Like: %v", err)
			}
			if changed != c.wantChanged {
				t.Fatalf("changed = %v, want %v", changed, c.wantChanged)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestTLikeModelUnlikeReportsChangeFromAffectedRows(t *testing.T) {
	const del = "DELETE FROM `t_like` WHERE target_type = \\? AND target_id = \\? AND user_id = \\?"

	cases := []struct {
		name        string
		affected    int64
		wantChanged bool
	}{
		{"确实取消", 1, true},
		{"本来没赞", 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			likeModel, mock := newTestLikeModel(t)
			mock.ExpectExec(del).WithArgs(1, 10, "u1").WillReturnResult(sqlmock.NewResult(0, c.affected))

			changed, err := likeModel.Unlike(context.Background(), 1, 10, "u1")
			if err != nil {
				t.Fatalf("Unlike: %v", err)
			}
			if changed != c.wantChanged {
				t.Fatalf("changed = %v, want %v", changed, c.wantChanged)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestTLikeModelLikedMapsHitTargets(t *testing.T) {
	likeModel, mock := newTestLikeModel(t)
	mock.ExpectQuery("SELECT \\* FROM `t_like` WHERE user_id = \\? AND target_type = \\? AND target_id IN").
		WillReturnRows(sqlmock.NewRows([]string{"id", "target_type", "target_id", "user_id"}).
			AddRow(1, 1, 10, "u1").
			AddRow(2, 1, 12, "u1"))

	got, err := likeModel.Liked(context.Background(), 1, "u1", []int64{10, 11, 12})
	if err != nil {
		t.Fatalf("Liked: %v", err)
	}
	if !got[10] || got[11] || !got[12] {
		t.Fatalf("Liked() = %v, want {10:true, 12:true}", got)
	}

	// 未登录不查库：再来一次也不会有新的期望被消耗
	got, err = likeModel.Liked(context.Background(), 1, "", []int64{10})
	if err != nil || len(got) != 0 {
		t.Fatalf("Liked() = %v, %v; want empty", got, err)
	}
}

func TestTLikeModelCountsGroupsByTarget(t *testing.T) {
	likeModel, mock := newTestLikeModel(t)
	mock.ExpectQuery("SELECT target_id, COUNT\\(\\*\\) AS total FROM `t_like` WHERE target_type = \\? AND target_id IN").
		WillReturnRows(sqlmock.NewRows([]string{"target_id", "total"}).AddRow(10, 3).AddRow(12, 1))

	got, err := likeModel.Counts(context.Background(), 1, []int64{10, 11, 12})
	if err != nil {
		t.Fatalf("Counts: %v", err)
	}
	if got[10] != 3 || got[12] != 1 || got[11] != 0 {
		t.Fatalf("Counts() = %v, want {10:3, 12:1}", got)
	}
}
