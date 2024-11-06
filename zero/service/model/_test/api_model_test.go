package _test_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/gormlogger"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/dbx"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
)

const dsn = "root:mysql7914@(127.0.0.1:3306)/blog-veweiyi?charset=utf8mb4&parseTime=True&loc=Local"

var db *gorm.DB

func init() {
	var err error
	// 连接数据库
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "tb_",
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
		Logger: logger.New(
			gormlogger.NewGormWriter(gormlogger.AddSkip(1)),
			logger.Config{
				SlowThreshold:             500 * time.Millisecond, // 慢 SQL 阈值，超过会提前结束
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false, // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  true,  // 彩色打印
				ParameterizedQueries:      false, // 使用参数化查询 (true时，会将参数值替换为?)
			},
		),
	})
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}
	log.Println("mysql connection done")
}

func Test_ClearTable(t *testing.T) {
	// 清空表
	//db.Exec("truncate table t_operation_log")

	//dbutil.CleanTable(db, "t_menu")
	//dbutil.CleanTable(db, "t_role_menu")

	dbx.CleanTable(db, "t_api")
	dbx.CleanTable(db, "t_role_api")

	//dbutil.CleanTable(db, "t_user_role")
}

func Test_TOperationLog(t *testing.T) {

	OperationLogModel := model.NewTOperationLogModel(db, nil)
	ctx := context.Background()

	data := &model.TOperationLog{
		UserId:         0,
		Nickname:       "test_nickname",
		IpAddress:      "test_ip_address",
		IpSource:       "test_ip_source",
		OptModule:      "test_opt_module",
		OptDesc:        "test_opt_desc",
		RequestUrl:     "test_request_url",
		RequestMethod:  "test_request_method",
		RequestHeader:  "test_request_header",
		RequestData:    "test_request_data",
		ResponseData:   "test_response_data",
		ResponseStatus: 200,
		Cost:           "",
	}

	batch, err := OperationLogModel.Deletes(ctx, "1=1")
	assert.Equal(t, nil, err)
	t.Log(batch)

	insert, err := OperationLogModel.Insert(ctx, data)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1), insert)
	t.Log(jsonconv.AnyToJsonIndent(data))

	data.Nickname = "test_nickname_update"
	data.IpAddress = ""
	data.ResponseStatus = 0
	update, err := OperationLogModel.Update(ctx, data)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1), update)

	one, err := OperationLogModel.FindOne(ctx, data.Id)
	assert.Equal(t, nil, err)
	t.Log(jsonconv.AnyToJsonIndent(one))

	data.Nickname = "test_nickname_save"
	data.IpAddress = ""
	data.ResponseStatus = 0
	save, err := OperationLogModel.Update(ctx, data)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1), save)

	first, err := OperationLogModel.First(ctx, "id = ?", data.Id)
	assert.Equal(t, nil, err)
	t.Log(jsonconv.AnyToJsonIndent(first))
}
