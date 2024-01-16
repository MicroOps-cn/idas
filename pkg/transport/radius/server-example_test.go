package radius_test

import (
	"fmt"
	"log"
	"testing"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

var (
	ServerUsername = "tim"
	ServerPassword = "12345"
)

type SecretSource struct {
}

func (s SecretSource) RADIUSSecret(r *radius.Request) ([]byte, error) {
	attr, err := radius.ParseAttributes(r.GetRaw()[20:])
	if err != nil {
		return nil, err
	}
	id, ok := attr.Lookup(rfc2865.NASIdentifier_Type)
	if !ok {
		return nil, fmt.Errorf("authentication failed")
	}
	if string(id) == "abcd" {
		return []byte("secret"), nil
	}
	if string(id) == "abcd1" {
		return []byte("secret1"), nil
	}
	return []byte("xxxx"), nil
}

func Example_packetServer() {
	handler := func(w radius.ResponseWriter, r *radius.Request) {
		username := rfc2865.UserName_GetString(r.Packet)
		password := rfc2865.UserPassword_GetString(r.Packet)

		var code radius.Code
		if username == ServerUsername && password == ServerPassword {
			code = radius.CodeAccessAccept
		} else {
			code = radius.CodeAccessReject
		}
		log.Printf("Writing %v to %v", code, r.RemoteAddr)
		w.Write(r.Response(code))
	}

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: &SecretSource{},
	}

	log.Printf("Starting server on :1812")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func TestExample_packetServer(t *testing.T) {

	Example_packetServer()
}
