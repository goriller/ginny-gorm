module github.com/goriller/ginny-gorm.git

go 1.16

require (
	github.com/google/wire v0.5.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/viper v1.12.0
	go.uber.org/zap v1.21.0
	gorm.io/driver/clickhouse v0.4.2
	gorm.io/driver/mysql v1.3.6
	gorm.io/driver/postgres v1.3.9
	gorm.io/driver/sqlite v1.3.6
	gorm.io/driver/sqlserver v1.3.2
	gorm.io/gorm v1.23.8
	moul.io/zapgorm2 v1.1.3
)
