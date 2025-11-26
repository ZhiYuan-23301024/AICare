package config

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DatabaseConfig 包含所有初始化数据库所需字段（可直接在 main 中构造并传入）
type DatabaseConfig struct {
	// 可直接传入完整 DSN（优先），否则会用下面的字段拼装 DSN
	DSN string `mapstructure:"dsn"`

	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`  // eg: "disable","require"
	TimeZone string `mapstructure:"timezone"` // eg: "UTC"

	// 连接池与日志配置
	MaxOpenConns        int    `mapstructure:"max_open_conns"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetimeSecs int    `mapstructure:"conn_max_lifetime_seconds"`
	LogLevel            string `mapstructure:"log_level"` // "silent","info","warn","error"
}

// buildDSN 根据字段构造 Postgres DSN（如果已设置 DSN 则直接返回）
func (c *DatabaseConfig) buildDSN() string {
	if strings.TrimSpace(c.DSN) != "" {
		return c.DSN
	}

	// 默认值处理
	host := c.Host
	if host == "" {
		host = "localhost"
	}
	port := c.Port
	if port == 0 {
		port = 5432
	}
	user := c.User
	if user == "" {
		user = "postgres"
	}
	dbname := c.DBName
	if dbname == "" {
		dbname = "postgres"
	}
	sslmode := c.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}
	timezone := c.TimeZone
	if timezone == "" {
		timezone = "UTC"
	}

	// Postgres DSN 格式： host=... user=... password=... dbname=... port=... sslmode=... TimeZone=...
	// 注意：如果 Password 包含空格或特殊字符，建议调用者传入完整 DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		host, user, c.Password, dbname, port, sslmode, timezone,
	)
	return dsn
}

// InitDatabase 用于初始化并返回 *gorm.DB
// 调用者只需构造 DatabaseConfig（包含密码/端口/数据库名等）并传入即可。
// 返回的 *gorm.DB 在进程退出时，需要通过 db.DB() 取出 *sql.DB 并 Close()。
func InitDatabase(cfg DatabaseConfig) (*gorm.DB, error) {
	dsn := strings.TrimSpace(cfg.buildDSN())
	if dsn == "" {
		return nil, fmt.Errorf("empty dsn")
	}

	// 选择 gorm 日志级别
	var lvl gormlogger.LogLevel
	switch strings.ToLower(cfg.LogLevel) {
	case "info":
		lvl = gormlogger.Info
	case "warn", "warning":
		lvl = gormlogger.Warn
	case "error":
		lvl = gormlogger.Error
	default:
		lvl = gormlogger.Silent
	}
	gormLogger := gormlogger.Default.LogMode(lvl)

	// 打开连接（gorm 会管理底层连接池）
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("open database failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB failed: %w", err)
	}

	// 应用连接池配置（若配置值 <=0 则保留默认）
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetimeSecs > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSecs) * time.Second)
	}

	// 启动时做一次 Ping 校验，超时后关闭并返回错误
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}
