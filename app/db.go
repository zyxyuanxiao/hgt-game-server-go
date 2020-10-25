package app

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"strconv"
	"time"
	"xorm.io/core"
)

var DB *xorm.EngineGroup

// mysql配置
type MysqlConn struct {
	Ip			string
	Username	string
	Password 	string
}

func LoadDB() {
	// ReadFile 函数会读取文件的全部内容，返回一个字节数组
	var err error
	// 数据库连接信息
	connections := MysqlConfig.handleMysqlConn(append([]MysqlConn{MysqlConfig.Master}, MysqlConfig.Slaves...))
	DB, err = xorm.NewEngineGroup("mysql", connections)
	if err != nil {
		fmt.Printf("数据库连接错误，%v", err)
	}
	//设置日志显示
	DB.ShowSQL(true)
	DB.SetLogLevel(core.LOG_DEBUG)
	//设置连接池DB
	DB.SetMaxOpenConns(3)
	DB.SetMaxIdleConns(1)
	DB.SetConnMaxLifetime(12 * time.Hour)
	//测试连接
	DB.Ping()
}

// 处理数据库链接
func (mc *MysqlConfigJson) handleMysqlConn(conn []MysqlConn) []string {
	var connections []string
	temp := ""
	for _, c := range conn {
		temp = c.Username + ":" + c.Password + "@" + mc.Protocol +
			"(" + c.Ip + ":" + strconv.FormatInt(mc.Port, 10) + ")/" +
			mc.Database + "?charset=" + mc.Charset + ";"
		connections = append(connections, temp)
	}

	return connections
}
