package gorm

import (
	"database/sql"
	"fmt"
	"ipfast_server/pkg/util/log"
	syslog "log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql" // Add this line to import the missing package
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

/*
全局数据库操作对象结构体

	DataBaseUrl string 数据库连接地址 Sqlite 为文件地址，Mysql 为连接地址
	Instance *gorm.DB 数据库操作对象
	SQLInstance *sql.DB 数据库连接,提供了一些基本的数据库操作方法，如 Query, Exec 等
	func Ping() error 检测数据库连接是否有效
	func Close() error 关闭数据库连接
	func SetLoglevel(levelstr string) error 设置数据库日志级别
	func BatchUpdate(updateDatas []interface{}, fields []string) error 批量更新数据
*/
type DataBaseInstance struct {
	DataBaseUrl string
	Instance    *gorm.DB
	SQLInstance *sql.DB
}

/*
全局Gorm数据库操作对象
*/
var MasterDb *DataBaseInstance

/*
根据配置确定数据库引擎

	param:
		config: 数据库配置
	return:
		dialector: 数据库引擎
*/
func DetermineDatabaseEngine() (dialector gorm.Dialector) {
	dataType := viper.GetString("database.type")
	if dataType == "mysql" {
		MasterDb.DataBaseUrl = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			viper.GetString("database.user"),
			viper.GetString("database.password"),
			viper.GetString("database.host"),
			viper.GetString("database.dbname"),
		)
		dialector = mysql.Open(MasterDb.DataBaseUrl)
	} else if dataType == "sqlite" {
		MasterDb.DataBaseUrl = viper.GetString("database.dbname")
		dialector = sqlite.Open(MasterDb.DataBaseUrl)
	}
	log.Debug("[ConnectDB]:%s", MasterDb.DataBaseUrl)
	return
}

/*
根据数据库引擎初始化数据库连接

	param:
		dialector: 数据库引擎
	return:
		newDb: 数据库操作对象
		err: 可能的错误
*/
func ConnectByDialector(dialector gorm.Dialector) (newDb *gorm.DB, err error) {
	return gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",   // 表名前缀，`User` 的表名将是 `t_users`
			SingularTable: true, // 使用单数表名，禁用表名复数化
		},
	})
}

func init() {
	MasterDb = &DataBaseInstance{}
}

/*
根据配置初始化数据库连接

	return:
		error: 可能的错误
*/
func Setup() (err error) {
	var newDb *gorm.DB
	// 根据配置确定数据库引擎
	newDb, err = ConnectByDialector(DetermineDatabaseEngine())
	if err != nil {
		return
	}
	// 关闭可能的旧数据库连接
	MasterDb.Close()
	// 更新数据库连接
	MasterDb.Instance = newDb
	MasterDb.SQLInstance, err = MasterDb.Instance.DB()

	if err == nil {
		// 连接池中最大空闲连接数
		MasterDb.SQLInstance.SetMaxIdleConns(5)
		// 最大打开连接数
		MasterDb.SQLInstance.SetMaxOpenConns(300)
		// 设置连接的最大生命周期。
		// MasterDb.SQLInstance.SetConnMaxLifetime(time.Hour * time.Duration(1))
	}
	// 设置批量插入的数量
	MasterDb.Instance.CreateBatchSize = 10000
	// 设置日志级别
	err = MasterDb.SetLoglevel(viper.GetString("database.log"))
	if err != nil {
		return
	}
	// 检测数据库连接是否有效
	err = MasterDb.Ping()
	if err != nil {
		return
	}
	// 连接成功
	log.Info("[数据库连接成功:%s ", MasterDb.DataBaseUrl)
	return
}

/*
检测数据库连接是否有效

	return:
		error: 可能的错误
*/
func (DBI *DataBaseInstance) Ping() error {
	if DBI.SQLInstance != nil {
		return DBI.SQLInstance.Ping()
	}
	return fmt.Errorf("sqlDB is invalid")
}

// CloneTable 创建一个新表，结构与现有表相同
func (DBI *DataBaseInstance) CloneTable(srcTableName, newTableName string) error {
	// 构建 SQL 语句以克隆表结构
	sqlStmt := fmt.Sprintf("CREATE TABLE %s LIKE %s;", newTableName, srcTableName)
	// 执行 SQL 语句
	return DBI.Instance.Exec(sqlStmt).Error
}

// RenameTable 重命名表
func (DBI *DataBaseInstance) RenameTable(oldTableName, newTableName string) error {
	// 构建 SQL 语句以重命名表
	sqlStmt := fmt.Sprintf("RENAME TABLE %s TO %s;", oldTableName, newTableName)
	// 执行 SQL 语句
	return DBI.Instance.Exec(sqlStmt).Error
}

func (DBI *DataBaseInstance) Model(value interface{}) *gorm.DB {
	return DBI.Instance.Model(value)
}

// DropTable 删除表
func (DBI *DataBaseInstance) DropTable(tableName string) error {
	// 检查表是否存在
	exists := DBI.Instance.Migrator().HasTable(tableName)
	if !exists {
		// 如果表不存在，返回一个错误或直接返回nil表示没有需要删除的表
		return nil
	}
	// 构建 SQL 语句以删除表
	sqlStmt := fmt.Sprintf("DROP TABLE %s;", tableName)
	// 执行 SQL 语句
	return DBI.Instance.Exec(sqlStmt).Error
}

// 自动迁移表 AutoMigrate 方法会检查 传入 结构体，并在数据库中创建或更新一个名为结构体名称 的表。如果表不存在，它会被创建。如果表已经存在，它的结构会被更新以匹配 传入 结构体
func (DBI *DataBaseInstance) AutoMigrate(dst ...interface{}) error {
	return DBI.Instance.AutoMigrate(dst...)
}

// 自动迁移表 AutoMigrate 方法会检查 传入 结构体，并在数据库中创建或更新一个名为结构体名称 的表。如果表不存在，它会被创建。如果表已经存在，它的结构会被更新以匹配 传入 结构体
func (DBI *DataBaseInstance) Order(value interface{}) *gorm.DB {
	return DBI.Instance.Order(value)
}

/*
检测数据库连接是否有效

	return:
		error: 可能的错误
*/
func (DBI *DataBaseInstance) Clauses(conds ...clause.Expression) (tx *gorm.DB) {
	return DBI.Instance.Clauses(conds...)
}

/*
向数据库里插入数据

	return:
		tx: 数据库操作对象
*/
func (DBI *DataBaseInstance) Create(value interface{}) (tx *gorm.DB) {
	return DBI.Instance.Create(value)
}

/*
设置操作表

	return:
		tx: 数据库操作对象
*/
func (DBI *DataBaseInstance) Table(name string, args ...interface{}) (tx *gorm.DB) {
	return DBI.Instance.Table(name, args...)
}

/*
设置操作表

	return:
		tx: 数据库操作对象
*/
func (DBI *DataBaseInstance) Preload(query string, args ...interface{}) (tx *gorm.DB) {
	return DBI.Instance.Preload(query, args...)
}

/*
设置查询数据库条件

	return:
		tx: 数据库操作对象
*/
func (DBI *DataBaseInstance) Where(query interface{}, args ...interface{}) (tx *gorm.DB) {
	return DBI.Instance.Where(query, args...)
}

/*
查询数据库

	return:
		tx: 数据库操作对象
*/
func (DBI *DataBaseInstance) Select(query interface{}, args ...interface{}) (tx *gorm.DB) {
	return DBI.Instance.Select(query, args...)
}

/*
查询数据库

	return:
		tx: 数据库操作对象
*/
func (DBI *DataBaseInstance) Find(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return DBI.Instance.Find(dest, conds)
}

/*
关闭数据库连接
*/
func (DBI *DataBaseInstance) Close() (err error) {
	log.Debug("DB Close")
	if DBI.SQLInstance != nil {
		err = DBI.SQLInstance.Close()
	}
	return
}

/*
设置数据库日志级别

	param:
		level: 日志级别
		"silent": logger.Silent,
		"error":  logger.Error,
		"warn":   logger.Warn,
		"info":   logger.Info,
*/
func (DBI *DataBaseInstance) SetLoglevel(levelstr string) error {
	levelstr = strings.ToLower(levelstr)
	levelMap := map[string]logger.LogLevel{
		"silent": logger.Silent,
		"error":  logger.Error,
		"warn":   logger.Warn,
		"info":   logger.Info,
	}
	levelLower := strings.ToLower(levelstr)
	logLevel, ok := levelMap[levelLower]
	if !ok {
		return fmt.Errorf("invalid log level: %s", levelstr)
	}
	newLogger := logger.New(
		syslog.New(os.Stdout, "\r\n", syslog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second * 5, // 慢 SQL 阈值
			LogLevel:      logLevel,        // Log level
			Colorful:      true,            // 禁用彩色打印
		},
	)
	DBI.Instance.Logger = newLogger
	return nil
}

/*
批量更新数据

	param:
		updateDatas: 需要更新的数据列表
		fields: 需要更新的字段
	return:
		error: 可能的错误
*/
func (DBI *DataBaseInstance) BatchUpdate(updateDatas []interface{}, fields []string) error {
	// 创建一个新的事务
	tx := DBI.Instance.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查事务开始是否成功
	if tx.Error != nil {
		return tx.Error
	}

	// 批量更新服务器
	for _, data := range updateDatas {
		err := tx.Select(fields).Updates(data).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	return tx.Commit().Error
}

func (DBI *DataBaseInstance) UpdateColumn(column string, value interface{}) *gorm.DB {
	return DBI.Instance.UpdateColumn(column, value)
}

func (DBI *DataBaseInstance) Begin(opts ...*sql.TxOptions) *gorm.DB {
	return DBI.Instance.Begin(opts...)
}
