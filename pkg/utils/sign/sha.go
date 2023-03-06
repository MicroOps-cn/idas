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

package sign

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
)

func SumSha1Hmac(secret, payload string) string {
	key := []byte(secret)
	hash := hmac.New(sha1.New, key)
	hash.Write([]byte(payload))
	signedBytes := hash.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}

func SumSha245Hmac(secret string, payload ...string) string {
	key := []byte(secret)
	hash := hmac.New(sha256.New, key)
	for _, p := range payload {
		hash.Write([]byte(p))
	}
	signedBytes := hash.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}
