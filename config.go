package orm

import (
	"fmt"
	"net/url"
	"strings"

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
	Type string `json:"type" mapstructure:"type"`
	// mysql://localhost:3306/dbname[?username=value1&password=value2&paramN=valueN]
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

	// parse dsn
	o.WDB, err = o.parseUrl(o.WDB)
	if err != nil {
		return nil, errors.Wrap(err, "parse dsn error")
	}

	rdb := []string{}
	for _, v := range o.RDBs {
		u, err := o.parseUrl(v)
		if err != nil {
			return nil, errors.Wrap(err, "parse dsn error")
		}
		rdb = append(rdb, u)
	}
	o.RDBs = rdb

	return o, nil
}

func (c *Config) parseUrl(dsn string) (string, error) {
	if dsn == "" {
		return "", fmt.Errorf("undefined wdb dsn")
	}
	u, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}
	q := u.Query()
	wdbUser := q.Get("username")
	wdbPass := q.Get("password")

	q.Del("username")
	q.Del("password")

	// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", wdbUser, wdbPass, u.Host, strings.TrimPrefix(u.Path, "/"), q.Encode()), nil
}
