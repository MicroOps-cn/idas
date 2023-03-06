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

package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/MicroOps-cn/idas/pkg/errors"
)

func (x AppMeta_GrantType) MarshalJSON() ([]byte, error) {
	var ret []string
	for val := range AppMeta_GrantType_name {
		if val != 0 && int32(x)&val == val {
			ret = append(ret, strconv.Itoa(int(val)))
		}
	}

	return []byte("[" + strings.Join(ret, ", ") + "]"), nil
}

func NewGrantType(gs ...AppMeta_GrantType) AppMeta_GrantType {
	var ret AppMeta_GrantType
	for _, g := range gs {
		ret |= g
	}
	return ret
}

//func (x AppMeta_GrantMode) MarshalJSON() ([]byte, error) {
//	return []byte(`"` + x.String() + `"`), nil
//}
//
//func (x AppMeta_Status) MarshalJSON() ([]byte, error) {
//	return []byte(`"` + x.String() + `"`), nil
//}
//
//func (x UserMeta_UserStatus) MarshalJSON() ([]byte, error) {
//	return []byte(`"` + x.String() + `"`), nil
//}

func (x *UserMeta_UserStatus) UnmarshalJSON(bytes []byte) error {
	if strings.HasPrefix(string(bytes), `"`) {
		var name string
		err := json.Unmarshal(bytes, &name)
		if err != nil {
			return err
		}
		s, ok := UserMeta_UserStatus_value[name]
		if !ok {
			return errors.ParameterError(fmt.Sprintf("unknown status: %s", name))
		}
		*x = UserMeta_UserStatus(s)
	} else {
		var val int32
		err := json.Unmarshal(bytes, &val)
		if err != nil {
			return err
		}

		if _, ok := UserMeta_UserStatus_name[val]; !ok {
			return errors.ParameterError(fmt.Sprintf("unknown status: %d", val))
		}
		*x = UserMeta_UserStatus(val)
	}

	return nil
}

func (x UserMeta_UserStatus) Is(s ...UserMeta_UserStatus) bool {
	for _, status := range s {
		if x&status != status {
			return false
		}
	}
	return true
}

func (x UserMeta_UserStatus) IsAnyOne(s ...UserMeta_UserStatus) bool {
	for _, status := range s {
		if x&status == status {
			return true
		}
	}
	return false
}

const (
	UserMetaStatusAll UserMeta_UserStatus = -1
)
