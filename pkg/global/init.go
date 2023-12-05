/*
 Copyright © 2023 MicroOps-cn.

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

package global

import (
	"regexp"

	"github.com/asaskevich/govalidator"
)

var rxHashLikeUUID = regexp.MustCompile("^[0-9a-f]{32}$")

func init() {
	govalidator.TagMap["uuid"] = func(str string) bool {
		switch len(str) {
		case 36:
			return govalidator.IsUUID(str)
		case 32:
			return rxHashLikeUUID.MatchString(str)
		default:
			return false
		}
	}
}
