package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
	Addr     Address

	Nickname string
	Birthday time.Time
	AboutMe  string
}

//领域对象和数据库对象存在差异，
//并不是严格映射，主要取决于服务的业务内容
// 比如说这里的address 结构体，在dao层中的实现是 单独一个json字段,或者另外一张表
//单对于domain来说，只关心json或者表中的数据是否按照要求给出带有
//Province string,Regin string的实例，并提供Address结构体给service使用

// type User struct {
// 	Id        int64 `gorm:"primaryKey"` //原生语言字面量用反引号
// 	Email     string
// 	Pawssword string
// 	Ctime     int64
// 	Utime     int64
// 	//json
// 	Addr string
// }

// //或者单独是一张表
// type Address struct{
// 	Uid int  `gorm:"primaryKey"`
// }
type Address struct {
	Province string
	Regin    string
}
