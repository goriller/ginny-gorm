package orm

import (
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

func newLogger(logger *zap.Logger) logger.Interface {
	l := zapgorm2.New(logger)
	l.SetAsDefault()
	return l
}
