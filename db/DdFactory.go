package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLDB() (*sql.DB, error) {
	// 设置数据库连接参数
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	password := os.Getenv("MYSQL_PASSWORD")
	user := os.Getenv("MYSQL_USER")

	db, err := sql.Open("mysql", ""+user+":"+password+"@tcp("+host+":"+port+")/simpleDouyin")
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败")
	}

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败")
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

func CloseDB() error {
	err := dB.Close()
	if err != nil {
		return err
	}

	fmt.Println("MySQL数据库连接关闭")

	return nil
}
