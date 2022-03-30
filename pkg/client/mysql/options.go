package mysql

import (
	"time"

	"github.com/gogo/protobuf/types"

	"idas/config"
)

func NewMySQLOptions() *config.MySQLOptions {
	return &config.MySQLOptions{
		Charset:               "utf8",
		Collation:             "utf8_general_ci",
		MaxIdleConnections:    2,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: types.DurationProto(30 * time.Second),
		TablePrefix:           "t_",
	}
}
