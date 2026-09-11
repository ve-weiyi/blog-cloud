package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TUserModel = (*defaultTUserModel)(nil)

type (
	// 接口定义
	TUserModel interface {
		TableName() string
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TUserModel)
		// 插入
		Insert(ctx context.Context, in *TUser) (rows int64, err error)
		InsertBatch(ctx context.Context, in ...*TUser) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TUser) (rows int64, err error)
		UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存或更新（更新零值）
		Save(ctx context.Context, in *TUser) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TUser, err error)
		FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TUser, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TUser, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TUser, total int64, err error)
		// add extra method in here
		FindOneByEmail(ctx context.Context, email string) (out *TUser, err error)
		FindOneByMobile(ctx context.Context, mobile string) (out *TUser, err error)
		FindOneByUserId(ctx context.Context, user_id string) (out *TUser, err error)
		FindOneByUsername(ctx context.Context, username string) (out *TUser, err error)
	}

	// 表字段定义
	TUser struct {
		Id           int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:ID" json:"id"`                                                          // ID
		UserId       string     `gorm:"column:user_id;type:varchar(64);not null;uniqueIndex:uk_uid,priority:1;default:'';comment:用户id" json:"user_id"`                     // 用户id
		Username     string     `gorm:"column:username;type:varchar(50);not null;uniqueIndex:uk_username,priority:1;default:'';comment:用户名" json:"username"`               // 用户名
		Password     string     `gorm:"column:password;type:varchar(255);not null;default:'';comment:加密后的密码（bcrypt）" json:"password"`                                      // 加密后的密码（bcrypt）
		Nickname     string     `gorm:"column:nickname;type:varchar(100);not null;default:'';comment:用户昵称" json:"nickname"`                                                // 用户昵称
		Avatar       string     `gorm:"column:avatar;type:varchar(500);not null;default:'';comment:用户头像URL" json:"avatar"`                                                 // 用户头像URL
		Email        *string    `gorm:"column:email;type:varchar(100);uniqueIndex:uk_email,priority:1;comment:邮箱" json:"email"`                                            // 邮箱
		Mobile       *string    `gorm:"column:mobile;type:varchar(20);uniqueIndex:uk_mobile,priority:1;comment:手机号" json:"mobile"`                                         // 手机号
		Info         string     `gorm:"column:info;type:varchar(1024);not null;default:'';comment:用户信息" json:"info"`                                                       // 用户信息
		Status       int64      `gorm:"column:status;type:bigint;not null;index:idx_status,priority:1;default:1;comment:账号状态：0-禁用 1-正常 2-冻结" json:"status"`                // 账号状态：0-禁用 1-正常 2-冻结
		RegisterType string     `gorm:"column:register_type;type:varchar(64);not null;default:'';comment:注册方式" json:"register_type"`                                       // 注册方式
		IpAddress    string     `gorm:"column:ip_address;type:varchar(255);not null;default:'';comment:注册ip" json:"ip_address"`                                            // 注册ip
		IpSource     string     `gorm:"column:ip_source;type:varchar(255);not null;default:'';comment:注册ip 源" json:"ip_source"`                                            // 注册ip 源
		CreatedAt    time.Time  `gorm:"column:created_at;type:datetime;not null;index:idx_created_at,priority:1;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
		UpdatedAt    time.Time  `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                 // 更新时间
		DeletedAt    *time.Time `gorm:"column:deleted_at;type:datetime;index:idx_deleted_at,priority:1;comment:删除时间，软删除" json:"deleted_at"`                                // 删除时间，软删除
	}

	// 接口实现
	defaultTUserModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTUserModel(db *gorm.DB) TUserModel {
	return &defaultTUserModel{
		DbEngin: db,
		table:   "`t_user`",
	}
}

func (m *defaultTUserModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTUserModel) WithTx(tx *gorm.DB) (out TUserModel) {
	return NewTUserModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTUserModel) Insert(ctx context.Context, in *TUser) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTUserModel) InsertBatch(ctx context.Context, in ...*TUser) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTUserModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TUser{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTUserModel) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TUser{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTUserModel) Update(ctx context.Context, in *TUser) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTUserModel) UpdateFields(ctx context.Context, feilds map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(feilds)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存或更新（更新零值）
func (m *defaultTUserModel) Save(ctx context.Context, in *TUser) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTUserModel) FindById(ctx context.Context, id int64) (out *TUser, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTUserModel) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TUser, err error) {
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
func (m *defaultTUserModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TUser, err error) {
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
func (m *defaultTUserModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTUserModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TUser, total int64, err error) {
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
func (m *defaultTUserModel) FindOneByEmail(ctx context.Context, email string) (out *TUser, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`email` = ?", email).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
func (m *defaultTUserModel) FindOneByMobile(ctx context.Context, mobile string) (out *TUser, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`mobile` = ?", mobile).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
func (m *defaultTUserModel) FindOneByUserId(ctx context.Context, user_id string) (out *TUser, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`user_id` = ?", user_id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
func (m *defaultTUserModel) FindOneByUsername(ctx context.Context, username string) (out *TUser, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`username` = ?", username).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
