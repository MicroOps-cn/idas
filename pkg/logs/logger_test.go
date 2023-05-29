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

package logs

import (
	"bytes"
	"testing"

	log "github.com/MicroOps-cn/fuck/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/require"
)

func TestRegisterLogger(t *testing.T) {
	log.RegisterLogFormat(FormatIDAS, NewIdasLogger)
	buf := bytes.NewBuffer(nil)
	l := log.New(log.WithWriter(buf), log.WithConfig(log.MustNewConfig("info", string(FormatIDAS))))
	level.Error(l).Log("msg", "test message", WrapKeyName("Name"), "Test")
	const matchExpr = `(?m)^[-.\d:TZ]+ \[error] [-\w]+ \S+ - test message - \n\[Name]:\s+Test`
	require.Regexp(t, matchExpr, buf.String())
}
