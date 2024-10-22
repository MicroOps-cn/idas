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
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"

	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/MicroOps-cn/idas/pkg/common"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/golang-jwt/jwt/v4"
)

type JWTIssuer interface {
	SignedString(ctx context.Context, claims Claims) (string, error)
	ParseWithClaims(tokenString string, claims jwt.Claims) (*jwt.Token, error)
	GetPublicKey() crypto.PublicKey
}

type JWTConfig struct {
	PrivateKey any
	PublicKey  crypto.PublicKey
	Algorithm  jwt.SigningMethod
	Id         string
}

func (j *JWTConfig) GetPublicKey() crypto.PublicKey {
	return j.PublicKey
}

func (j *JWTConfig) ParseWithClaims(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		switch token.Method.(type) {
		case *jwt.SigningMethodRSA, *jwt.SigningMethodECDSA:
			if j.PublicKey == nil {
				return nil, fmt.Errorf("public key is nil")
			}

			return j.PublicKey, nil
		case *jwt.SigningMethodHMAC:
			return j.PrivateKey, nil
		default:
			return "", fmt.Errorf("invalid algorithm: %s", j.Algorithm)
		}
	})
}

func (j *JWTConfig) SignedString(ctx context.Context, claims Claims) (string, error) {
	claims.SetIssuer(ctx, j.Id)
	return jwt.NewWithClaims(j.Algorithm, claims).SignedString(j.PrivateKey)
}

func (j *JWTConfig) UnmarshalJSON(bytes []byte) (err error) {
	type plain struct {
		Secret     string `json:"secret"`
		PrivateKey string `json:"private_key"`
		Algorithm  string `json:"algorithm"`
	}
	var c plain
	if err = json.Unmarshal(bytes, &c); err != nil {
		return err
	}
	if c.Algorithm == "" {
		if c.PrivateKey != "" {
			block, _ := pem.Decode([]byte(c.PrivateKey))
			if _, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
				c.Algorithm = "RS256"
			} else if _, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
				c.Algorithm = "ES256"
			} else if privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
				switch privKey.(type) {
				case *ecdsa.PrivateKey:
					c.Algorithm = "ES256"
				case *rsa.PrivateKey:
					c.Algorithm = "RS256"
				default:
					return fmt.Errorf("unsupported private key type: %T", privKey)
				}
			}
		} else if c.Secret != "" {
			c.Algorithm = "HS256"
		}
	}
	issuer, err := NewJWTConfig("", c.Algorithm, w.DefaultString(c.PrivateKey, c.Secret))
	if err != nil {
		return err
	}
	*j = *issuer
	return nil
}
func privateKeyToPEM(pk any) (string, error) {
	asn1Bytes, err := x509.MarshalPKCS8PrivateKey(pk)
	if err != nil {
		return "", err
	}
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: asn1Bytes,
	}

	// 将PEM块编码为PEM格式
	pemBytes := pem.EncodeToMemory(pemBlock)

	return string(pemBytes), nil
}

//func (j JWTConfig) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
//	return j.MarshalJSON()
//}

func (j JWTConfig) MarshalJSON() ([]byte, error) {
	type plain struct {
		PrivateKey string `json:"private_key"`
		Algorithm  string `json:"algorithm"`
	}
	switch j.Algorithm.(type) {
	case *jwt.SigningMethodHMAC:
		return json.Marshal(plain{
			PrivateKey: string(j.PrivateKey.([]byte)),
			Algorithm:  j.Algorithm.Alg(),
		})
	case *jwt.SigningMethodRSA, *jwt.SigningMethodECDSA:
		pemStr, err := privateKeyToPEM(j.PrivateKey)
		if err != nil {
			return nil, err
		}
		return json.Marshal(plain{
			PrivateKey: pemStr,
			Algorithm:  j.Algorithm.Alg(),
		})
	default:
		return nil, errors.New("unsupported algorithm")
	}
}

func NewRandomKey(method string) (string, error) {
	switch method {
	case "HS256", "HS384", "HS512":
		priv := make([]byte, 2048)
		_, err := rand.Read(priv)
		if err != nil {
			return "", err
		}
		return string(priv), nil
	case "RS256", "RS384", "RS512":
		priv, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return "", err
		}
		privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
		if err != nil {
			return "", fmt.Errorf("unable to marshal private key: %v", err)
		}
		return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})), nil
	case "ES256", "ES384", "ES512":
		priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return "", err
		}
		privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
		if err != nil {
			return "", fmt.Errorf("unable to marshal private key: %v", err)
		}
		return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})), nil
	default:
		return "", fmt.Errorf("invalid algorithm: %s", method)
	}
}

func NewRandomRSAJWTConfig() (*JWTConfig, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &JWTConfig{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		Algorithm:  jwt.SigningMethodRS256,
	}, nil
}

func NewJWTConfigBySecret(secret string) (*JWTConfig, error) {
	return &JWTConfig{PrivateKey: []byte(secret), Algorithm: jwt.SigningMethodHS256}, nil
}

var tokenAlgorithmMap = map[string]jwt.SigningMethod{
	"HS256": jwt.SigningMethodHS256,
	"HS384": jwt.SigningMethodHS384,
	"HS512": jwt.SigningMethodHS512,
	"RS256": jwt.SigningMethodRS256,
	"RS384": jwt.SigningMethodRS384,
	"RS512": jwt.SigningMethodRS512,
	"ES256": jwt.SigningMethodES256,
	"ES384": jwt.SigningMethodES384,
	"ES512": jwt.SigningMethodES512,
}

func NewJWTConfig(issuerId, method, privateKey string) (*JWTConfig, error) {
	var jwtConfig JWTConfig
	if alg, ok := tokenAlgorithmMap[method]; ok {
		jwtConfig.Algorithm = alg
	} else {
		return nil, fmt.Errorf("invalid algorithm: %s", method)
	}
	jwtConfig.Id = issuerId
	switch method {
	case "HS256", "HS384", "HS512":
		jwtConfig.PrivateKey = []byte(privateKey)
	case "RS256", "RS384", "RS512":
		privk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
		if err != nil {
			return nil, fmt.Errorf("failed to load rsa private key: %s", err)
		}
		jwtConfig.PublicKey = privk.Public()
		jwtConfig.PrivateKey = privk
	case "ES256", "ES384", "ES512":
		privk, err := jwt.ParseECPrivateKeyFromPEM([]byte(privateKey))
		if err != nil {
			return nil, fmt.Errorf("failed to load ecdsa private key: %s", err)
		}
		jwtConfig.PublicKey = privk.Public()
		jwtConfig.PrivateKey = privk
		switch privk.Curve.Params().BitSize {
		case 256:
			jwtConfig.Algorithm = jwt.SigningMethodES256
		case 384:
			jwtConfig.Algorithm = jwt.SigningMethodES384
		case 512:
			jwtConfig.Algorithm = jwt.SigningMethodES512
		default:
			return nil, fmt.Errorf("invalid ecdsa curve size: %d", privk.Curve.Params().BitSize)
		}
	}
	return &jwtConfig, nil
}

func NewJWTIssuer(issuerId string, method, privateKey string) (JWTIssuer, error) {
	return NewJWTConfig(issuerId, method, privateKey)
}

type Claims interface {
	jwt.Claims
	SetIssuer(context.Context, string)
}

func ParseWithClaims(tokenString string, claims jwt.Claims, issuerFunc func(token *jwt.Token) (JWTIssuer, error)) (*jwt.Token, error) {
	token, parts, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, err
	}
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid token")
	}
	issuer, err := issuerFunc(token)
	if err != nil {
		return nil, err
	}
	return issuer.ParseWithClaims(tokenString, claims)
}

type StandardClaims jwt.StandardClaims

func (c StandardClaims) Valid() error {
	return jwt.StandardClaims(c).Valid()
}

func (c *StandardClaims) SetIssuer(ctx context.Context, issuerID string) {
	if len(issuerID) > 0 {
		c.Issuer = w.M(common.GetURL(ctx, common.WithAPI("v1", "oauth"), common.WithParam("client_id", issuerID)))
	} else {
		c.Issuer = w.M(common.GetURL(ctx, common.WithAPI("v1", "oauth")))
	}
}
