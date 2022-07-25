package sign

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
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
	r, s := new(big.Int), new(big.Int)
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
		return "", fmt.Errorf("私钥格式异常")
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
