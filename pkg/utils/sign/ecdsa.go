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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
)

func ECDSAVerify(pub1 string, pub2 string, payload string, sig string) bool {
	var r, s big.Int
	var ok bool
	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int),
		Y:     new(big.Int),
	}
	data := []byte(payload)
	sigPair := strings.Split(sig, ":")
	if len(sigPair) != 2 {
		return false
	}
	if _, ok = r.SetString(sigPair[0], 62); !ok {
		return false
	}
	if _, ok = s.SetString(sigPair[1], 62); !ok {
		return false
	}
	if _, ok = publicKey.X.SetString(pub1, 62); !ok {
		return false
	}
	if _, ok = publicKey.Y.SetString(pub2, 62); !ok {
		return false
	}
	return ecdsa.Verify(publicKey, data, &r, &s)
}

func ECDSASign(priv string, payload string) (hash string, err error) {
	var r, s *big.Int
	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int),
			Y:     new(big.Int),
		},
		D: new(big.Int),
	}
	var ok bool
	if _, ok = privateKey.D.SetString(priv, 62); !ok {
		return "", errors.New("private key format exception")
	}
	r, s, err = ecdsa.Sign(rand.Reader, privateKey, []byte(payload))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s",
		r.Text(62),
		s.Text(62),
	), nil
}

func GenerateECDSAKeyPair() (pub1, pub2, private string, err error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", "", err
	}
	return privateKey.X.Text(62), privateKey.Y.Text(62), privateKey.D.Text(62), nil
}
