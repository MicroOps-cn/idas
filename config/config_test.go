package config

import (
	"bytes"
	"fmt"
	"idas/pkg/client/gorm"
	"idas/pkg/utils/wrapper"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/stretchr/testify/require"

	"idas/pkg/logs"
	"idas/pkg/utils/httputil"
)

var conf = `
storage:
  user:
  - mysql:
      maxIdleConnections: 2
      maxOpenConnections: 100
      maxConnectionLifeTime: 30s
      charset: utf8
      collation: utf8_general_ci
      tablePrefix: t_idas_
`

func TestUnmarshalConfig(t *testing.T) {
	logger := logs.New(logs.MustNewConfig("info", "json"))
	err := safeCfg.ReloadConfigFromYamlReader(logger, NewConverter("./", bytes.NewReader([]byte(conf))))
	require.NoError(t, err)
	require.Equal(t, safeCfg.C.Storage.User[0].GetMysql().TablePrefix, "t_idas_")
}

func TestMarshalConfig(t *testing.T) {
	mysqlOptions := gorm.NewMySQLOptions()
	mysqlOptions.TablePrefix = "t_xsadfa9i83"
	c := Config{
		Storage: &Storages{
			User: []*Storage{{
				Source: &Storage_Mysql{
					Mysql: mysqlOptions,
				},
			}},
		},
	}

	marshaler := jsonpb.Marshaler{
		Indent:   "    ",
		OrigName: true,
	}
	buf := bytes.NewBuffer(nil)
	err := marshaler.Marshal(buf, &c)
	require.NoError(t, err)

	t.Log(buf.String())
	logger := logs.New(logs.MustNewConfig("info", "json"))
	err = safeCfg.ReloadConfigFromYamlReader(logger, NewConverter("./", bytes.NewReader([]byte(conf))))
	require.NoError(t, err)
	t.Log(safeCfg.C.Storage)
	require.Equal(t, c.Storage.User[0].GetMysql().TablePrefix, "t_xsadfa9i83")
	require.NoError(t, err)
}

func TestValues(t *testing.T) {
	ints := wrapper.Must[[]time.Duration](httputil.NewValue("10s,2m,60s,30a0m").Split().Durations())
	fmt.Println(ints)
}
