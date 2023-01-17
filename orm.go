package orm

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"

	"github.com/google/wire"
	"github.com/goriller/ginny-gorm/dialector"
	"github.com/goriller/ginny-util/graceful"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Provider = wire.NewSet(ConfigProvider, New)

type ORM struct {
	writeDB *gorm.DB
	readDBs []*gorm.DB
	logger  *zap.Logger
}

// New
func New(ctx context.Context, conf *Config, logger *zap.Logger) (*ORM, error) {
	writeDB, err := newDB(ctx, conf.WDB, conf)
	if err != nil {
		return nil, err
	}
	// RDB多个
	rDBLen := len(conf.RDBs)
	readDBs := make([]*gorm.DB, 0, rDBLen)
	for i := 0; i < rDBLen; i++ {
		readDB, err := newDB(ctx, conf.RDBs[i], conf)
		if err != nil {
			return nil, err
		}
		readDBs = append(readDBs, readDB)
	}
	db := &ORM{
		writeDB: writeDB,
		readDBs: readDBs,
		logger:  logger,
	}
	// graceful
	graceful.AddCloser(func(ctx context.Context) error {
		return db.Close()
	})

	return db, nil
}

// newDB
func newDB(ctx context.Context, dsn string, conf *Config) (*gorm.DB, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "faild parse dsn")
	}
	switch conf.Type {
	case "sqllite":
		conf.dialector = dialector.NewSqllite(u)
		break
	case "mysql":
		conf.dialector = dialector.NewMysql(u)
		break
	case "postgres":
	case "postgresql":
		conf.dialector = dialector.NewPostgres(u)
		break
	case "sqlserver":
		conf.dialector = dialector.NewSqlserver(u)
		break
	case "clickhouse":
		conf.dialector = dialector.NewClickhouse(u)
		break
	default:
		return nil, errors.Wrap(err, "not support")
	}
	db, err := gorm.Open(conf.dialector, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database")
	}
	return db, nil
}

// RDB 随机返回一个读库
func (m *ORM) RDB() *gorm.DB {
	return m.readDBs[rand.Intn(len(m.readDBs))]
}

// WDB 返回唯一写库
func (m *ORM) WDB() *gorm.DB {
	return m.writeDB
}

// Close 关闭所有读写连接池，停止keepalive保活协程。该函数应当很少使用到
func (m *ORM) Close() error {
	db, err := m.writeDB.DB()
	if err != nil {
		return err
	}
	if err := db.Close(); err != nil {
		m.logger.Error("close write db error", zap.Error(err))
		return err
	}
	for i := 0; i < len(m.readDBs); i++ {
		rdb, err := m.readDBs[i].DB()
		if err != nil {
			return err
		}
		if err := rdb.Close(); err != nil {
			m.logger.Error("close db read pool error", zap.Error(err))
			return err
		}
	}
	return nil
}
