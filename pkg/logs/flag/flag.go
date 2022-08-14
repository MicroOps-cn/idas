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

package flag

import (
	"github.com/spf13/pflag"

	"github.com/MicroOps-cn/idas/pkg/logs"
)

// LevelFlagName is the canonical flag name to configure the allowed log level
// within Prometheus projects.
const LevelFlagName = "log.level"

// LevelFlagHelp is the help description for the log.level flag.
const LevelFlagHelp = "Only log messages with the given severity or above. One of: [debug, info, warn, error]"

// FormatFlagName is the canonical flag name to configure the log format
// within Prometheus projects.
const FormatFlagName = "log.format"

// FormatFlagHelp is the help description for the log.format flag.
const FormatFlagHelp = "Output format of log messages. One of: [logfmt, json, idas]"

// AddFlags adds the flags used by this package to the Kingpin application.
// To use the default Kingpin application, call AddFlags(kingpin.CommandLine)
func AddFlags(set *pflag.FlagSet, config *logs.Config) {
	config.Level = new(logs.AllowedLevel)
	set.StringVar((*string)(config.Level), LevelFlagName, string(logs.LevelInfo), LevelFlagHelp)
	config.Format = new(logs.AllowedFormat)
	set.StringVar((*string)(config.Format), FormatFlagName, string(logs.FormatIdas), FormatFlagHelp)
}
