package svc

import (
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/model"
)

// AutoMigrate 按模型补全表结构；生产环境以 SQL 迁移脚本为准。
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.TUser{},
		&model.TUserOauth{},
		&model.TRole{},
		&model.TApi{},
		&model.TMenu{},
		&model.TUserRole{},
		&model.TRoleApi{},
		&model.TRoleMenu{},
		&model.TArticle{},
		&model.TCategory{},
		&model.TTag{},
		&model.TArticleTag{},
		&model.TChat{},
		&model.TComment{},
		&model.TMessage{},
		&model.TConfig{},
		&model.TAlbum{},
		&model.TPhoto{},
		&model.TFriend{},
		&model.TTalk{},
		&model.TPage{},
		&model.TGuest{},
		&model.TNotificationEvent{},
		&model.TNotification{},
		&model.TNotificationDelivery{},
		&model.TNotificationPreference{},
		&model.TAnnouncement{},
		&model.TDailyStats{},
		&model.TVisitLog{},
		&model.TLoginLog{},
		&model.TOperationLog{},
		&model.TUploadLog{},
	)
}
