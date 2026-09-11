package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TCommentModel = (*defaultTCommentModel)(nil)

type (
	// 接口定义
	TCommentModel interface {
		TableName() string
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TCommentModel)
		// 插入
		Insert(ctx context.Context, in *TComment) (rows int64, err error)
		InsertBatch(ctx context.Context, in ...*TComment) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TComment) (rows int64, err error)
		UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存或更新（更新零值）
		Save(ctx context.Context, in *TComment) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TComment, err error)
		FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TComment, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TComment, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TComment, total int64, err error)
		// add extra method in here
	}

	// 表字段定义
	TComment struct {
		Id             int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键" json:"id"`                                           // 主键
		UserId         string    `gorm:"column:user_id;type:varchar(64);not null;index:fk_comment_user,priority:1;default:'';comment:评论用户id" json:"user_id"` // 评论用户id
		DeviceId       string    `gorm:"column:device_id;type:varchar(64);not null;default:'';comment:设备id" json:"device_id"`                                // 设备id
		TopicId        int64     `gorm:"column:topic_id;type:bigint;not null;default:0;comment:主题id" json:"topic_id"`                                        // 主题id
		ParentId       int64     `gorm:"column:parent_id;type:bigint;not null;index:fk_comment_parent,priority:1;default:0;comment:父评论id" json:"parent_id"`  // 父评论id
		ReplyId        int64     `gorm:"column:reply_id;type:bigint;not null;default:0;comment:回复评论id" json:"reply_id"`                                      // 回复评论id
		ReplyUserId    string    `gorm:"column:reply_user_id;type:varchar(255);not null;default:'';comment:评论回复用户id" json:"reply_user_id"`                   // 评论回复用户id
		CommentContent string    `gorm:"column:comment_content;type:text;not null;default:'';comment:评论内容" json:"comment_content"`                           // 评论内容
		Type           int64     `gorm:"column:type;type:bigint;not null;default:0;comment:评论类型 1.文章 2.友链 3.说说" json:"type"`                                 // 评论类型 1.文章 2.友链 3.说说
		Status         int64     `gorm:"column:status;type:bigint;not null;default:0;comment:状态 0.正常 1.已编辑 2.已删除" json:"status"`                             // 状态 0.正常 1.已编辑 2.已删除
		LikeCount      int64     `gorm:"column:like_count;type:bigint;not null;default:0;comment:评论点赞数量" json:"like_count"`                                  // 评论点赞数量
		CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                  // 创建时间
		UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                  // 更新时间
	}

	// 接口实现
	defaultTCommentModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTCommentModel(db *gorm.DB) TCommentModel {
	return &defaultTCommentModel{
		DbEngin: db,
		table:   "`t_comment`",
	}
}

func (m *defaultTCommentModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTCommentModel) WithTx(tx *gorm.DB) (out TCommentModel) {
	return NewTCommentModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTCommentModel) Insert(ctx context.Context, in *TComment) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTCommentModel) InsertBatch(ctx context.Context, in ...*TComment) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTCommentModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TComment{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTCommentModel) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TComment{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTCommentModel) Update(ctx context.Context, in *TComment) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTCommentModel) UpdateFields(ctx context.Context, feilds map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(feilds)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存或更新（更新零值）
func (m *defaultTCommentModel) Save(ctx context.Context, in *TComment) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTCommentModel) FindById(ctx context.Context, id int64) (out *TComment, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTCommentModel) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TComment, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	// 如果有排序参数
	if len(sorts) != 0 {
		db = db.Order(sorts)
	}

	err = db.First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询列表
func (m *defaultTCommentModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TComment, err error) {
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

// 查询总数
func (m *defaultTCommentModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 分页查询记录
func (m *defaultTCommentModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TComment, total int64, err error) {
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

	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
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
		return nil, 0, err
	}

	return list, total, nil
}

// add extra method in here
