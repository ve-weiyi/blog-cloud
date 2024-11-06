package dbx

import (
	"fmt"

	"gorm.io/gorm"
)

func CleanTable(db *gorm.DB, tableName string) (err error) {
	// 清空表的数据
	err = db.Exec(fmt.Sprintf("DELETE FROM `%v`", tableName)).Error
	if err != nil {
		return err
	}
	// 重置 AUTO_INCREMENT 值为 1
	err = db.Exec(fmt.Sprintf("ALTER TABLE `%v` AUTO_INCREMENT = 1", tableName)).Error
	if err != nil {
		return err
	}
	return nil
}
