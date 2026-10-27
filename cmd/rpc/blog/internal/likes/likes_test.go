package likes

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/model"
)

// fakeLikeModel 只实现 Apply 用到的方法，其余由内嵌接口兜底。
type fakeLikeModel struct {
	model.TLikeModel
	rows map[string]bool
}

func (f *fakeLikeModel) key(targetType, targetID int64, userID string) string {
	return string(rune(targetType)) + ":" + string(rune(targetID)) + ":" + userID
}

func (f *fakeLikeModel) FindOneByTargetTypeTargetIdUserId(_ context.Context, targetType, targetID int64, userID string) (*model.TLike, error) {
	if f.rows[f.key(targetType, targetID, userID)] {
		return &model.TLike{}, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeLikeModel) Insert(_ context.Context, in *model.TLike) (int64, error) {
	f.rows[f.key(in.TargetType, in.TargetId, in.UserId)] = true
	return 1, nil
}

func (f *fakeLikeModel) DeleteBatch(_ context.Context, _ string, args ...interface{}) (int64, error) {
	f.rows[f.key(args[0].(int64), args[1].(int64), args[2].(string))] = false
	return 1, nil
}

func TestApplyConverges(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name      string
		existing  bool
		wantLike  bool
		wantDelta int64
	}{
		{"未赞时点赞", false, true, 1},
		{"重复点赞不再计数", true, true, 0},
		{"已赞时取消", true, false, -1},
		{"重复取消不再计数", false, false, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := &fakeLikeModel{rows: map[string]bool{}}
			if c.existing {
				m.rows[m.key(1, 10, "u1")] = true
			}

			delta, err := Apply(ctx, m, 1, 10, "u1", c.wantLike)
			if err != nil {
				t.Fatalf("Apply: %v", err)
			}
			if delta != c.wantDelta {
				t.Fatalf("delta = %d, want %d", delta, c.wantDelta)
			}
			if got := m.rows[m.key(1, 10, "u1")]; got != c.wantLike {
				t.Fatalf("最终状态 = %v, want %v", got, c.wantLike)
			}
		})
	}

	// 查询失败不能被当成「不存在」，否则会把重复点赞计成新增
	broken := &brokenLikeModel{}
	if _, err := Apply(ctx, broken, 1, 10, "u1", true); !errors.Is(err, errBoom) {
		t.Fatalf("查询失败时 err = %v, want %v", err, errBoom)
	}
}

var errBoom = errors.New("boom")

type brokenLikeModel struct {
	model.TLikeModel
}

func (b *brokenLikeModel) FindOneByTargetTypeTargetIdUserId(context.Context, int64, int64, string) (*model.TLike, error) {
	return nil, errBoom
}
