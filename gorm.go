package sns

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var dbVariable *gorm.DB

// AutoMigrate 创建表结构
func AutoMigrate() {
	db := BeginDB()
	db.AutoMigrate(&CodeInfo{})
	db.DBCommit()
}

// InitMysqlWithDB 初始化MySQL
func InitMysqlWithDB(db *gorm.DB) {
	dbVariable = db
}

// InitMysql 初始化MySQL
func InitMysql(mysqlHost string, debug bool) error {
	db, err := gorm.Open("mysql", mysqlHost)
	if err != nil {
		return err
	}
	if debug {
		db = db.Debug()
	}
	dbVariable = db
	return nil
}

// CloseMysql 关闭MySQL链接
func CloseMysql() {
	if dbVariable != nil {
		dbVariable.Close()
	}
}

// BaseDB 基础DB
type BaseDB struct {
	*gorm.DB
	skip bool
}

// DBCommit 提交事务
func (b *BaseDB) DBCommit() {
	if b.skip {
		return
	}
	b.skip = true
	b.DB.Commit()
}

// DBRollback 回滚事物
func (b *BaseDB) DBRollback() {
	if b.skip {
		return
	}
	b.skip = true
	b.DB.Rollback()
}

// BeginDB 开始事物
func BeginDB() *BaseDB {
	db := dbVariable.Begin()
	return &BaseDB{DB: db}
}
