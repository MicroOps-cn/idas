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

package httputil

import (
	"net/url"

	w "github.com/MicroOps-cn/fuck/wrapper"
)

type Map[KT comparable, VT any] map[KT]VT

func (m Map[KT, VT]) String() string {
	return w.JSONStringer(m).String()
}

type URL url.URL

func (u URL) String() string {
	return (*url.URL)(&u).String()
}

func (u *URL) Set(s string) error {
	ou, err := url.Parse(s)
	if err != nil {
		return err
	}
	*u = URL(*ou)
	return nil
}

func (u *URL) Type() string {
	return "URL"
}
