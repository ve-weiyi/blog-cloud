package model

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ TChatMessageModel = (*defaultTChatMessageModel)(nil)

type (
	// 接口定义
	TChatMessageModel interface {
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out TChatMessageModel)
		// 插入
		Insert(ctx context.Context, in *TChatMessage) (rows int64, err error)
		Inserts(ctx context.Context, in ...*TChatMessage) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TChatMessage) (rows int64, err error)
		Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存
		Save(ctx context.Context, in *TChatMessage) (rows int64, err error)
		// 查询
		FindOne(ctx context.Context, id int64) (out *TChatMessage, err error)
		First(ctx context.Context, conditions string, args ...interface{}) (out *TChatMessage, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TChatMessage, err error)
		FindList(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TChatMessage, err error)
		// add extra method in here
	}

	// 表字段定义
	TChatMessage struct {
		Id          int64     `json:"id" gorm:"column:id" `                       // 主键
		UserId      string    `json:"user_id" gorm:"column:user_id" `             // 用户id
		DeviceId    string    `json:"device_id" gorm:"column:device_id" `         // 设备id
		TopicId     string    `json:"topic_id" gorm:"column:topic_id" `           // 主题id,表示一个群会话
		ReplyMsgId  string    `json:"reply_msg_id" gorm:"column:reply_msg_id" `   // 回复消息id，at消息
		ReplyUserId string    `json:"reply_user_id" gorm:"column:reply_user_id" ` // 回复用户列表,at用户列表
		ChatContent string    `json:"chat_content" gorm:"column:chat_content" `   // 聊天内容
		IpAddress   string    `json:"ip_address" gorm:"column:ip_address" `       // 用户ip 127.0.0.1
		IpSource    string    `json:"ip_source" gorm:"column:ip_source" `         // 用户地址 广东省深圳市
		Type        string    `json:"type" gorm:"column:type" `                   // 类型:chatgpt chatroom
		Status      int64     `json:"status" gorm:"column:status" `               // 状态:0正常 1编辑 2撤回 3删除
		CreatedAt   time.Time `json:"created_at" gorm:"column:created_at" `       // 创建时间
		UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at" `       // 更新时间
	}

	// 接口实现
	defaultTChatMessageModel struct {
		DbEngin    *gorm.DB
		CacheEngin *redis.Client
		table      string
	}
)

func NewTChatMessageModel(db *gorm.DB, cache *redis.Client) TChatMessageModel {
	return &defaultTChatMessageModel{
		DbEngin:    db,
		CacheEngin: cache,
		table:      "`t_chat_message`",
	}
}

func (m *defaultTChatMessageModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTChatMessageModel) WithTransaction(tx *gorm.DB) (out TChatMessageModel) {
	return NewTChatMessageModel(tx, m.CacheEngin)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTChatMessageModel) Insert(ctx context.Context, in *TChatMessage) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTChatMessageModel) Inserts(ctx context.Context, in ...*TChatMessage) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTChatMessageModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TChatMessage{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTChatMessageModel) Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TChatMessage{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTChatMessageModel) Update(ctx context.Context, in *TChatMessage) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTChatMessageModel) Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(columns)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *defaultTChatMessageModel) Save(ctx context.Context, in *TChatMessage) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTChatMessageModel) FindOne(ctx context.Context, id int64) (out *TChatMessage, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTChatMessageModel) First(ctx context.Context, conditions string, args ...interface{}) (out *TChatMessage, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.First(&out).Error
	if err != nil {
		return nil, err
	}
	return out, err
}

// 查询总数
func (m *defaultTChatMessageModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Model(&TChatMessage{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 查询列表
func (m *defaultTChatMessageModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TChatMessage, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, err
}

// 分页查询记录
func (m *defaultTChatMessageModel) FindList(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TChatMessage, err error) {
	// 插入db
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有搜索条件
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	// 如果有排序参数
	if len(sorts) != 0 {
		db = db.Order(sorts)
	}

	// 如果有分页参数
	if page > 0 && size > 0 {
		limit := size
		offset := (page - 1) * limit
		db = db.Limit(limit).Offset(offset)
	}

	// 查询数据
	err = db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

// add extra method in here
