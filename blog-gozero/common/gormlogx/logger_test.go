package gormlogx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func TestLogger(t *testing.T) {
	dsn := "root:mysql7914@(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true"

	lg := New(
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 超过这个耗时的查询会打印日志
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false, // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false, // 彩色打印
			ParameterizedQueries:      false, // 使用参数化查询 (true时，会将参数值替换为?)
		},
	)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
		Logger: lg,
	})
	if err != nil {
		return
	}

	traceId, err := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	assert.Equal(t, nil, err)
	spanId, err := trace.SpanIDFromHex("00f067aa0ba902b7")
	assert.Equal(t, nil, err)
	state, err := trace.ParseTraceState("key1=value1,key2=value2")
	assert.Equal(t, nil, err)
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceId,
		SpanID:     spanId,
		TraceState: state,
		Remote:     true,
	})

	ctx := context.Background()
	ctx = trace.ContextWithRemoteSpanContext(ctx, sc)

	logx.WithContext(ctx).Info("我是Info日志")
	gormDB.Logger.Warn(ctx, "我是Info日志")

	var result []Tag
	// 实际项目需要传入来自go-zero启用trace的ctx，才能获取到traceId和spanID
	err = gormDB.WithContext(ctx).Table("tag").First(&result).Error
	assert.Equal(t, nil, err)
	t.Log(result)

	var tag Tag
	err = gormDB.WithContext(ctx).Table("tag").Where("id = 0").First(&tag).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

type Tag struct {
	ID      int    `gorm:"column:id"`
	TagName string `gorm:"column:tag_name"`
}
