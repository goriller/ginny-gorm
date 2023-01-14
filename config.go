package orm

import (
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ConfigProvider
var ConfigProvider = wire.NewSet(NewConfig)

// Config
type Config struct {
	// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	WDB  string   `json:"wdb" mapstructure:"wdb"`
	RDBs []string `json:"rdbs" mapstructure:"rdbs"`
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
	o.QueryFields = true
	o.Logger = newLogger(logger)

	if o.RDBs == nil || len(o.RDBs) == 0 {
		o.RDBs = []string{
			o.WDB,
		}
	}
	return o, nil
}
