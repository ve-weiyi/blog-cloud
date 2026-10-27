package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TOperationLogModel = (*defaultTOperationLogModel)(nil)

type (
	// 接口定义
	TOperationLogModel interface {
		Table[TOperationLog]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TOperationLogModel)
	}

	// 表字段定义
	TOperationLog struct {
		Id             int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键id" json:"id"`                                                        // 主键id
		UserId         string    `gorm:"column:user_id;type:varchar(64);not null;index:idx_user_id,priority:1;default:'';comment:用户id" json:"user_id"`                      // 用户id
		DeviceId       string    `gorm:"column:device_id;type:varchar(64);not null;default:'';comment:设备id" json:"device_id"`                                               // 设备id
		Module         string    `gorm:"column:module;type:varchar(50);not null;default:'';comment:操作模块" json:"module"`                                                     // 操作模块
		Description    string    `gorm:"column:description;type:varchar(500);not null;default:'';comment:操作描述" json:"description"`                                          // 操作描述
		RequestUri     string    `gorm:"column:request_uri;type:varchar(255);not null;default:'';comment:请求地址" json:"request_uri"`                                          // 请求地址
		RequestMethod  string    `gorm:"column:request_method;type:varchar(32);not null;default:'';comment:请求方式" json:"request_method"`                                     // 请求方式
		RequestData    string    `gorm:"column:request_data;type:varchar(4096);not null;default:'';comment:请求参数" json:"request_data"`                                       // 请求参数
		ResponseData   string    `gorm:"column:response_data;type:varchar(4096);not null;default:'';comment:返回数据" json:"response_data"`                                     // 返回数据
		ResponseStatus int64     `gorm:"column:response_status;type:smallint;not null;default:0;comment:响应状态码" json:"response_status"`                                      // 响应状态码
		Cost           int64     `gorm:"column:cost;type:int;not null;default:0;comment:耗时（ms）" json:"cost"`                                                                // 耗时（ms）
		CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;index:idx_created_at,priority:1;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
		UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                 // 更新时间
	}

	// 接口实现
	defaultTOperationLogModel struct {
		baseModel[TOperationLog]
	}
)

func NewTOperationLogModel(db *gorm.DB) TOperationLogModel {
	return &defaultTOperationLogModel{newBaseModel[TOperationLog](db, "`t_operation_log`")}
}

// 在事务中操作
func (m *defaultTOperationLogModel) WithTx(tx *gorm.DB) (out TOperationLogModel) {
	return &defaultTOperationLogModel{newBaseModel[TOperationLog](tx, m.table)}
}
