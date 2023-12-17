package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {

	//db = db.Debug()
	return db.AutoMigrate(&User{})
}
