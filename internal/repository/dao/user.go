package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var ErrDuplicateEmail = errors.New("DuplicateEmail")
var ErrRecordNotFind = gorm.ErrRecordNotFound

type UserDAO struct {
	db *gorm.DB
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {

	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error //处理mysql的error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 邮箱主键冲突
			return ErrDuplicateEmail
		}
	}
	return err //未定义错误
}

type User struct {
	//自增组件相邻 操作系统预读
	//建表主要看索引使用是否正确
	Id       int64  `gorm:"primaryKey autoIncrement"` //原生语言字面量用反引号
	Email    string `gorm:"unique"`
	Password string
	//UTC 0 毫秒数v
	Ctime int64
	Utime int64
	//json
	Addr     string
	Nickname string `gorm:"type=varchar(128)"`
	// YYYY-MM-DD
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`
}

// 或者单独是一张表
type Address struct {
	Uid int `gorm:"primaryKey"`
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err //未定义错误
}

func (dao *UserDAO) UpdateById(ctx context.Context, entity User) error {
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.Id).
		Updates(map[string]any{
			"utime":    time.Now().UnixMilli(),
			"nickname": entity.Nickname,
			"birthday": entity.Birthday,
			"about_me": entity.AboutMe,
		}).Error
}

func (dao *UserDAO) FindById(ctx context.Context, uid int64) (User, error) {
	var res User
	err := dao.db.WithContext(ctx).Where("id = ?", uid).First(&res).Error
	return res, err
}
