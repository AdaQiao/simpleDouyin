package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLDB() (*sql.DB, error) {
	// 设置数据库连接参数
	db, err := sql.Open("mysql", "root:tCzAhYFo@tcp(172.16.32.27:49721)/simpleDouyin")
	if err != nil {
		return nil, err
	}

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("MySQL数据库连接成功")

	return db, nil
}

var (
	dB *sql.DB
)

func InitDB() error {
	db, err := NewMySQLDB()
	if err != nil {
		return err
	}

	dB = db

	return nil
}

func GetDB() *sql.DB {
	return dB
}
