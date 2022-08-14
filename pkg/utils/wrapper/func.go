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

package w

func M[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Error[T any](_ T, err error) error {
	return err
}

func ToInterfaces[T any](objs []T) []interface{} {
	var newObjs []interface{}
	for _, obj := range objs {
		newObjs = append(newObjs, obj)
	}
	return newObjs
}

func P[T any](o T) *T {
	return &o
}
