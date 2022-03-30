package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-kit/log"
	"github.com/golang/protobuf/jsonpb"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v3"

	"idas/pkg/global"
)

var (
	configReloadSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: global.AppName,
		Name:      "config_last_reload_successful",
		Help:      "Blackbox exporter config loaded successfully.",
	})

	configReloadSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: global.AppName,
		Name:      "config_last_reload_success_timestamp_seconds",
		Help:      "Timestamp of the last successful configuration reload.",
	})
	safeCfg = newSafeConfig()
)

type safeConfig struct {
	sync.RWMutex
	C *Config
}

func newSafeConfig() *safeConfig {
	return &safeConfig{
		C: &Config{},
	}
}

func Get() *Config {
	return safeCfg.GetConfig()
}

func (sc *safeConfig) SetConfig(conf *Config) {
	sc.Lock()
	defer sc.Unlock()
	sc.C = conf
}

func (sc *safeConfig) GetConfig() *Config {
	sc.RLock()
	defer sc.RUnlock()
	return sc.C
}

type Converter map[string]interface{}

func (c Converter) ToJSON() io.Reader {
	jsonReader := bytes.NewBuffer(nil)
	if err := json.NewEncoder(jsonReader).Encode(c); err != nil {
		panic(err)
	}
	return jsonReader
}

func (sc *safeConfig) ReloadConfigFromYamlReader(logger log.Logger, yamlReader io.Reader) (err error) {
	defer func() {
		if err != nil {
			configReloadSuccess.Set(0)
		} else {
			configReloadSuccess.Set(1)
			configReloadSeconds.SetToCurrentTime()
		}
	}()
	cfgConvert := new(Converter)
	if err = yaml.NewDecoder(yamlReader).Decode(&cfgConvert); err != nil {
		return fmt.Errorf("error parse config file: %s", err)
	}
	return sc.ReloadConfigFromJSONReader(logger, cfgConvert.ToJSON())
}

func (sc *safeConfig) ReloadConfigFromJSONReader(logger log.Logger, reader io.Reader) (err error) {
	defer func() {
		if err != nil {
			configReloadSuccess.Set(0)
		} else {
			configReloadSuccess.Set(1)
			configReloadSeconds.SetToCurrentTime()
		}
	}()

	var c Config

	var unmarshaler jsonpb.Unmarshaler
	if err = unmarshaler.Unmarshal(reader, &c); err != nil {
		return fmt.Errorf("error unmarshal config: %s", err)
	} else if err = c.Init(logger); err != nil {
		return fmt.Errorf("error init config: %s", err)
	}
	sc.SetConfig(&c)
	return nil
}

func (sc *safeConfig) ReloadConfigFromFile(logger log.Logger, filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}
	ext := filepath.Ext(filename)
	if len(ext) > 1 {
		switch ext {
		case ".yml", ".yaml":
			return sc.ReloadConfigFromYamlReader(logger, r)
		case ".json":
			return sc.ReloadConfigFromJSONReader(logger, r)
		}
	}
	return nil
}

func ReloadConfigFromFile(logger log.Logger, filename string) error {
	return safeCfg.ReloadConfigFromFile(logger, filename)
}
