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

package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type JWTIssuer interface {
	SignedString(claims jwt.Claims) (string, error)
	ParseWithClaims(tokenString string, claims jwt.Claims) (*jwt.Token, error)
	GetPublicKey() *rsa.PublicKey
}

type JWTConfig struct {
	Secret     string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	Algorithm  jwt.SigningMethod
}

func (j *JWTConfig) GetPublicKey() *rsa.PublicKey {
	return j.PublicKey
}

func (j *JWTConfig) ParseWithClaims(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		switch token.Method.(type) {
		case *jwt.SigningMethodRSA:
			if j.PublicKey == nil {
				return nil, fmt.Errorf("public key is nil")
			}
			return j.PublicKey, nil
		case *jwt.SigningMethodHMAC:
			return []byte(j.Secret), nil
		default:
			return "", fmt.Errorf("invalid algorithm: %s", j.Algorithm)
		}
	})
}

func (j *JWTConfig) SignedString(claims jwt.Claims) (string, error) {
	switch j.Algorithm.(type) {
	case *jwt.SigningMethodRSA:
		return jwt.NewWithClaims(j.Algorithm, claims).SignedString(j.PrivateKey)
	case *jwt.SigningMethodHMAC:
		return jwt.NewWithClaims(j.Algorithm, claims).SignedString([]byte(j.Secret))
	default:
		return "", fmt.Errorf("invalid algorithm: %s", j.Algorithm)
	}
}

func (j *JWTConfig) UnmarshalJSON(bytes []byte) (err error) {
	fmt.Println(bytes)
	type plain struct {
		Secret     string `json:"secret"`
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
		Algorithm  string `json:"algorithm"`
	}
	var c plain
	if err = json.Unmarshal(bytes, &c); err != nil {
		return err
	}
	if len(c.PublicKey) > 0 && len(c.PrivateKey) > 0 {
		if j.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(c.PublicKey)); err != nil {
			return fmt.Errorf("failed to load public key: %s", err)
		}
		if j.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(c.PrivateKey)); err != nil {
			return fmt.Errorf("failed to load private key: %s", err)
		}
		if j.PublicKey.N.Cmp(j.PrivateKey.N) != 0 || j.PublicKey.E != j.PrivateKey.E {
			return fmt.Errorf("public key does not match private key")
		}

		switch c.Algorithm {
		case "RS256", "":
			j.Algorithm = jwt.SigningMethodRS256
		case "RS512":
			j.Algorithm = jwt.SigningMethodRS512
		case "RS384":
			j.Algorithm = jwt.SigningMethodRS384
		default:
			return fmt.Errorf("invalid algorithm: %s", c.Algorithm)
		}
	} else if len(c.Secret) > 0 {
		j.Secret = c.Secret
		switch c.Algorithm {
		case "HS256", "":
			j.Algorithm = jwt.SigningMethodHS256
		case "HS512":
			j.Algorithm = jwt.SigningMethodHS512
		case "HS384":
			j.Algorithm = jwt.SigningMethodHS384
		default:
			return fmt.Errorf("invalid algorithm: %s", c.Algorithm)
		}
	}
	return nil
}

func NewJWTConfig() (*JWTConfig, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	fmt.Println(x509.MarshalPKCS1PrivateKey(privateKey))
	fmt.Println(x509.MarshalPKCS1PublicKey(&privateKey.PublicKey))

	return &JWTConfig{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		Algorithm:  jwt.SigningMethodHS256,
	}, nil
}

func NewJWTConfigBySecret(secret string) (*JWTConfig, error) {
	return &JWTConfig{Secret: secret, Algorithm: jwt.SigningMethodHS256}, nil
}
