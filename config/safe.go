/*
 Copyright © 2022 MicroOps-cn.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v3"

	"github.com/MicroOps-cn/idas/pkg/global"
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

type Converter struct {
	io.Reader
	name string
}

func (c *Converter) Name() string {
	return c.name
}

func (c *Converter) UnmarshalYAML(value *yaml.Node) error {
	vals := make(map[string]interface{})
	err := value.Decode(&vals)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	c.Reader = buf
	return json.NewEncoder(buf).Encode(vals)
}

var _ yaml.Unmarshaler = &Converter{}

func NewConverter(name string, r io.Reader) *Converter {
	return &Converter{name: name, Reader: r}
}

func (sc *safeConfig) ReloadConfigFromYamlReader(logger log.Logger, yamlReader Reader) (err error) {
	defer func() {
		if err != nil {
			configReloadSuccess.Set(0)
		} else {
			configReloadSuccess.Set(1)
			configReloadSeconds.SetToCurrentTime()
		}
	}()
	cfgConvert := new(Converter)
	cfgConvert.name = yamlReader.Name()
	if err = yaml.NewDecoder(yamlReader).Decode(&cfgConvert); err != nil {
		return fmt.Errorf("error parse config file: %s", err)
	}
	return sc.ReloadConfigFromJSONReader(logger, cfgConvert)
}

type Reader interface {
	io.Reader
	Name() string
}

func (sc *safeConfig) ReloadConfigFromJSONReader(logger log.Logger, reader Reader) (err error) {
	defer func() {
		if err != nil {
			configReloadSuccess.Set(0)
		} else {
			configReloadSuccess.Set(1)
			configReloadSeconds.SetToCurrentTime()
		}
	}()

	c := Config{
		Global: NewGlobalOptions(),
	}

	var unmarshaler jsonpb.Unmarshaler
	if err = unmarshaler.Unmarshal(reader, &c); err != nil {
		return fmt.Errorf("error unmarshal config: %s", err)
	} else if err = c.Init(logger); err != nil {
		return fmt.Errorf("error init config: %s", err)
	}
	if c.GetWorkspace() == nil {
		if absPath, err := filepath.Abs(path.Dir(reader.Name())); err != nil {
			c.SetWorkspace(path.Dir(reader.Name()))
			level.Debug(logger).Log("msg", "set workspace", "workspace", path.Dir(reader.Name()))
		} else {
			c.SetWorkspace(absPath)
			level.Debug(logger).Log("msg", "set workspace", "workspace", absPath)
		}
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

func ReloadConfigFromYamlReader(logger log.Logger, r Reader) error {
	return safeCfg.ReloadConfigFromYamlReader(logger, r)
}
