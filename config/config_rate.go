/*
 Copyright © 2024 MicroOps-cn.

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
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/gogo/protobuf/jsonpb"
	"golang.org/x/time/rate"
)

type nopLimiter struct{}

func (n nopLimiter) String() string {
	return "unlimited"
}

func (n nopLimiter) Limit() rate.Limit {
	return 0
}

func (n nopLimiter) Burst() int {
	return 0
}

func (n nopLimiter) Allow() bool {
	return true
}

var nopAllower = &nopLimiter{}

type Allowers []Limiter

func (l Allowers) String() string {
	return strings.Join(w.Map(l, func(item Limiter) string {
		return fmt.Sprintf("%vrps/%d", item.Limit(), item.Burst())
	}), ",")
}

func (l Allowers) Allow() bool {
	for _, limiter := range l {
		if !limiter.Allow() {
			return false
		}
	}
	return true
}

type stdLimiterWrapper struct {
	*rate.Limiter
}

func NewLimiterWrapper(l *rate.Limiter) Allower {
	return &stdLimiterWrapper{
		Limiter: l,
	}
}

func (s stdLimiterWrapper) String() string {
	return fmt.Sprintf("%v/%d", s.Limit(), s.Burst())
}

func (x *Config) GetRateLimit(name string) Allower {
	var allower Allowers
	for _, limit := range x.Security.RateLimit {
		if limit.Name.Contains(name) {
			if name == "" {
				if limit.Allower == nopAllower {
					return nopAllower
				}
				return &stdLimiterWrapper{Limiter: rate.NewLimiter(limit.Allower.Limit(), limit.Allower.Burst())}
			}
			allower = append(allower, limit.Allower)
		}
	}
	if len(allower) == 0 {
		return nil
	}
	return allower
}

func (r RateLimit) MarshalJSON() ([]byte, error) {
	type plain struct {
		Name  w.OneOrMore[string]
		Limit json.RawMessage
		Burst int
	}
	tmp := plain{
		Name:  r.Name,
		Limit: json.RawMessage(r.Limit),
		Burst: int(r.Burst),
	}
	return json.Marshal(tmp)
}

type Limiter interface {
	Allow() bool
	Limit() rate.Limit
	Burst() int
}

type Allower interface {
	Allow() bool
	String() string
}

func (r *RateLimit) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, b []byte) error {
	return r.UnmarshalJSON(b)
}

func (r *RateLimit) UnmarshalJSON(data []byte) (err error) {
	type plain struct {
		Name  w.OneOrMore[string]
		Limit json.RawMessage
		limit float64
		Burst int
	}
	var tmp plain
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	r.Name = tmp.Name
	if len(r.Name) == 0 {
		r.Name = []string{""}
	}
	r.Limit = string(tmp.Limit)

	if tmp.limit, err = strconv.ParseFloat(string(tmp.Limit), 64); err != nil {
		var tmpStr string
		if err := json.Unmarshal(tmp.Limit, &tmpStr); err != nil {
			return err
		}
		switch tmpStr {
		case "false", "0":
			r.Allower = nopAllower
			return nil
		default:
			if tmp.limit, err = strconv.ParseFloat(tmpStr, 64); err != nil {
				return err
			}
		}
	}
	if tmp.limit <= 0 {
		tmp.limit = 10
	}
	if tmp.Burst <= 0 {
		if tmp.limit > 256 {
			tmp.Burst = int(tmp.limit)
		} else {
			tmp.Burst = int(tmp.limit * (10 - math.Log(tmp.limit)))
		}
	}
	r.Allower = rate.NewLimiter(rate.Limit(tmp.limit), tmp.Burst)

	return nil
}
