package orm

import (
	"fmt"

	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Provider = wire.NewSet(ConfigProvider, New)

func New(conf *Config, logger *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(conf.dialector, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database")
	}

	return db, nil
}
