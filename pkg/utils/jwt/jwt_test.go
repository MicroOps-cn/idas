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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func generateRandomRSAKeyPair(t *testing.T) ([]byte, []byte) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	require.NoError(t, err)
	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return privateKeyPem, publicKeyPem
}

func generateRandomECDSAKeyPair(t *testing.T) ([]byte, []byte) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	privateKeyBytes, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	require.NoError(t, err)
	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "ECDSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return privateKeyPem, publicKeyPem
}

type MapClaims jwt.MapClaims

func (m MapClaims) Valid() error {
	return jwt.MapClaims(m).Valid()
}

func (m *MapClaims) SetIssuer(s string) {
	(*m)["iss"] = s
}

func TestJWTIssuer(t *testing.T) {
	rsaPrivK, rsaPubK := generateRandomRSAKeyPair(t)
	ecdsaPrivK, ecdsaPubK := generateRandomECDSAKeyPair(t)
	type args struct {
		method     string
		publicKey  string
		privateKey string
	}
	tests := []struct {
		name    string
		args    args
		want    JWTIssuer
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				method:     "HS256",
				publicKey:  "",
				privateKey: "secret",
			},
		},
		{
			name: "test-rsa",
			args: args{
				method:     "RS384",
				publicKey:  string(rsaPubK),
				privateKey: string(rsaPrivK),
			},
		},
		{
			name: "test-rsa",
			args: args{
				method:     "RS384",
				publicKey:  "-----BEGIN CERTIFICATE-----\nMIIFMzCCBBugAwIBAgIIc1JRIdv6BykwDQYJKoZIhvcNAQELBQAwgZMxCzAJBgNV\nBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAwDgYDVQQHEwdCZWlqaW5nMSwwKgYD\nVQQKEyNCZWlqaW5nIFdpc2Vhc3kgVGVjaG5vbG9neSBDby4sIEx0ZDEMMAoGA1UE\nCxMDT3BzMSQwIgYDVQQDExtXaXNlYXN5IFRlc3QgV2ViIFNlcnZpY2UgQ0EwHhcN\nMjIwMTA1MDQwNzMzWhcNMzIwMTAzMDQwNzMzWjCBjDELMAkGA1UEBhMCQ04xEDAO\nBgNVBAgTB0JlaWppbmcxEDAOBgNVBAcTB0JlaWppbmcxLDAqBgNVBAoTI0JlaWpp\nbmcgV2lzZWFzeSBUZWNobm9sb2d5IENvLiwgTHRkMQwwCgYDVQQLEwNPcHMxHTAb\nBgNVBAMTFHJlc291cmNlLndhbmd0ZXN0LmNuMIIBIjANBgkqhkiG9w0BAQEFAAOC\nAQ8AMIIBCgKCAQEA73QFHaeNoNt0PRyv5BOD7D+U8uzlCiSMvo9BsIDnYE+hECTb\n9S2sO6Y+KF6YDCAiW4QqZEZTu+Y4NOHEgAuMRKPa+lTAnpzcT1ZNrY6LhmL4dRib\nQ4DRGv5+YepTBTwP/HwUVojupUBWqkYFuIsIy/7Vs27vqM7zQB5VjFUguDkEs9+m\nrXhqKW/R96JKdkKr44ahf2GA7cwGqZfxJw2fLqnxq1s89bLpHEPPhTGuRGxJ7WQM\njfv8x5+AmK2ZUs1bPT1ymXMa6fyFZMLpGyPv4m/v5BWt8fmoB6fzSKHvw60H4e1p\nwadpZd7EHW/Gf/BN1MzkBDcO2sfI9MNO98LuHwIDAQABo4IBjjCCAYowDgYDVR0P\nAQH/BAQDAgeAMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMB\nAf8EAjAAMB8GA1UdIwQYMBaAFGq5WSEvu+zYjReBAPFCjbGeXNN+MIG+BggrBgEF\nBQcBAQSBsTCBrjBUBggrBgEFBQcwAYZIaHR0cDovL29wcy53aXNlYXN5LmNvbS9h\ncGkvZ2F0ZXdheS9jZXJ0L29jc3AvV2lzZWFzeSBUZXN0IFdlYiBTZXJ2aWNlIENB\nMFYGCCsGAQUFBzAChkpodHRwOi8vb3BzLndpc2Vhc3kuY29tL2FwaS9nYXRld2F5\nL2NlcnQvY2EvV2lzZWFzeSBUZXN0IFdlYiBTZXJ2aWNlIENBLmNydDALBgNVHREE\nBDACggAwXAYDVR0fBFUwUzBRoE+gTYZLaHR0cDovL29wcy53aXNlYXN5LmNvbS9h\ncGkvZ2F0ZXdheS9jZXJ0L2NybC9XaXNlYXN5IFRlc3QgV2ViIFNlcnZpY2UgQ0Eu\nY3JsMA0GCSqGSIb3DQEBCwUAA4IBAQBP5zmgHJ6iqkhZgvwIa9ay+Fx1dwVNw965\nZ/89sAZkHAe+TFpwZ1xR8zqAHiXGQR8CbzBqB1W01TFJY8KnibXKm3JlNSlOK7x+\ncoxlavD5n8gY9qPVOOQftP4caRgTJmilN3uim7AqSPsx34NPE4ZccQ1D8svm4/RQ\niOOjxs+w4/IG9lopa1JBeBi5A11FT+9iKKlzxJbafWb7W7cW2O4ohNdFTaDf6lwc\nRCEuTI2EoJjtA04rl+wIay295k/DwojnhMb9A7O1FDngR1og44bOjO67F/MHRWAc\nmDu6U2x8phou17AHeH7OO/calDcEQOlDH4yG78428A6mXydjWGTt\n-----END CERTIFICATE-----",
				privateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEA73QFHaeNoNt0PRyv5BOD7D+U8uzlCiSMvo9BsIDnYE+hECTb\n9S2sO6Y+KF6YDCAiW4QqZEZTu+Y4NOHEgAuMRKPa+lTAnpzcT1ZNrY6LhmL4dRib\nQ4DRGv5+YepTBTwP/HwUVojupUBWqkYFuIsIy/7Vs27vqM7zQB5VjFUguDkEs9+m\nrXhqKW/R96JKdkKr44ahf2GA7cwGqZfxJw2fLqnxq1s89bLpHEPPhTGuRGxJ7WQM\njfv8x5+AmK2ZUs1bPT1ymXMa6fyFZMLpGyPv4m/v5BWt8fmoB6fzSKHvw60H4e1p\nwadpZd7EHW/Gf/BN1MzkBDcO2sfI9MNO98LuHwIDAQABAoIBAQCecXG1JppzduLa\nUTIdw8AGQign+hKv/HFY4mgAB7uSIf6cNReKi1cs/RqiEb2gQF8bmT+HrHVZnsNQ\nUpd4dquw+485F32BNqAcqympDupJ2RE4QjjymLlEmGM+HRQkIZMeaWf3vpHSrNjr\nwHummfEPMqdrHJveYlnY8nl+6xFEc8Z6bAV+sgDUs5VVx2Pfx2y1AjnAiPQmIE9Y\nqteoyVZKzHvA2NEy9CCHNG/plCYVesFkZ/S98mSyy7k0jB3WUPEK+gTcIN9yONyy\nwvwC8BB7Shox5RG1ul0b9zWIvF6n1cyLd5vYmiqoIryIk9upBI5JBH+eTCjJjUAT\nrVx94wJZAoGBAPm/Bc5egDBdWXR7PLcbw2j/nikpbLKEzLs5bx/aUpM17kBp4i3n\nHlI8fm1e8na8YHvf4jMPiX/8rby3GS+4sxpg1rj1psGAsFrg5JvDxNtSP9qhMlOr\n2jBO7rT/iZH6dKRqPT7PacXoQ7ksbVuZjb1R/fAakkRqb92dWJRPZ7TtAoGBAPVz\nA9eSfVH6rKQwRLgtgvi9xEonsXfjx2YRgmxLjMyanGzDR2Y6Q3TVUPxc/pg5nWIJ\npnelYRMY+Qa8Pq/xrjOjSqfSlAQNxo+wfm6+OburBuv9o9EguX5g7r4fnl1lKsjw\natAiIoLWV8gW+MUiQXxvpDPiBJKVMS14txyofTm7AoGBAJtYi8cDFx9+YU9H/Ms1\nFMayAXI/FyKv4h0vK4UXqzdwW2NruUmuMjka8dUcMxtSL32+FBiIuJGI3ZS+G4eI\njreAtu9Ttcc1Qf01WF3fVwrJTXizvfc3tT9JScgCD1NjA7zlbHUuVO/Kep2rGdbZ\nW8YAQ0Ffdc3imvSxk9Ck17A9AoGBAIh+bmuKJjDZovoncX+up3/mH+tRCYrvW2qy\nYAITPXhmnoiJTAJYcjzdh4zftiE3IQNs9GriyAoTwCBzvLShRMuoihKrsu5SLtKn\nRpgVJwvq/w1rXpckiKL0CrAl6y5q3REjSXL3GJQD2IsH403VT++AMiM8FGjjmJZ9\n4+6G8CSTAoGAfT8r8AYRuH52nfqHSXcA/9bmsnkXUZToZDT1nfws6QZv0japRU9I\nWviSllSKK/w3XzOqAz3Vw+XT5uDiMgeVznBurdFubORAFpZ97KwG2oIsUp7NpmaR\nMHgHBfxfhvmsxXkx8f6JuKh6DJHVLgo8bOZQjngYRYNDRIePrKM4kdg=\n-----END RSA PRIVATE KEY-----\n",
			},
		},
		{
			name: "test-ecdsa",
			args: args{
				method:     "ES512",
				publicKey:  string(ecdsaPubK),
				privateKey: string(ecdsaPrivK),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJWTIssuer("", tt.args.method, tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJWTIssuer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			claims := MapClaims{"test": "test"}
			signedString, err := got.SignedString(&claims)
			require.NoError(t, err)
			gotClaims, err := got.ParseWithClaims(signedString, jwt.MapClaims{})
			require.NoError(t, err)
			require.Equal(t, jwt.MapClaims(claims), gotClaims.Claims)
		})
	}
}
