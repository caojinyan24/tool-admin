package db

import (
	"context"
	"fmt"
	"time"
	"tooladmin/server/common/config"
	"tooladmin/server/common/logger"
	"tooladmin/server/scheduler/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(ctx context.Context, mysqlConfig config.MysqlConfig) error {
	logger.Info(ctx, "Initializing database connection")
	var err error

	// // 配置GORM日志级别
	// newLogger := gormLogger.New(
	// 	gormLogger.Config{
	// 		SlowThreshold: time.Second, // 慢SQL阈值
	// 		LogLevel:      gormLogger.LogLevel(mysqlConfig.LogLevel),
	// 		Colorful:      true,
	// 	},
	// )

	// 连接数据库
	DB, err = gorm.Open(mysql.Open(mysqlConfig.Master), &gorm.Config{
		// Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移表结构
	if err = migrateTables(); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqlDB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

// migrateTables 迁移表结构
func migrateTables() error {
	return DB.AutoMigrate(
		&model.Task{},
		&model.Log{},
		&model.Schedule{},
	)
}
