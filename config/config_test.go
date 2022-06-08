package config

import (
	"bytes"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/require"
	"idas/pkg/client/gorm"
	"testing"

	"idas/pkg/logs"
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
	require.Equal(t, safeCfg.C.Storage.User[0].GetSource().(*Storage_Mysql).Mysql.Options().TablePrefix, "t_idas_")
}

func TestMarshalConfig(t *testing.T) {
	mysqlOptions := gorm.NewMySQLOptions()
	mysqlOptions.TablePrefix = "t_xsadfa9i83"
	client := &gorm.MySQLClient{}
	client.SetOptions(mysqlOptions)
	c := Config{
		Storage: &Storages{
			User: []*Storage{{
				Source: &Storage_Mysql{
					Mysql: client,
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
	require.Equal(t, safeCfg.C.Storage.User[0].GetSource().(*Storage_Mysql).Mysql.Options().TablePrefix, "t_xsadfa9i83")
	require.NoError(t, err)
}
