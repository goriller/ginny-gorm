package orm

import (
	"net/url"

	"github.com/google/wire"
	"github.com/goriller/ginny-gorm.git/dialector"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ConfigProvider
var ConfigProvider = wire.NewSet(NewConfig)

// Config
type Config struct {
	Dsn string
	gorm.Config
	dialector gorm.Dialector
}

// NewConfig
func NewConfig(v *viper.Viper, logger *zap.Logger) (*Config, error) {
	var err error
	o := new(Config)
	if err = v.UnmarshalKey("gorm", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal gorm option error")
	}

	o.Logger = newLogger(logger)

	u, err := url.Parse(o.Dsn)
	if err != nil {
		return nil, errors.Wrap(err, "faild parse dsn")
	}
	switch u.Scheme {
	case "sqllite":
		o.dialector = dialector.NewSqllite(u)
		break
	case "mysql":
		o.dialector = dialector.NewMysql(u)
		break
	case "postgres":
	case "postgresql":
		o.dialector = dialector.NewPostgres(u)
		break
	case "sqlserver":
		o.dialector = dialector.NewSqlserver(u)
		break
	case "clickhouse":
		o.dialector = dialector.NewClickhouse(u)
		break
	default:
		return nil, errors.Wrap(err, "not support")
	}

	return o, nil
}
