package radius_test

import (
	"context"
	"log"
	"testing"
	"time"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

var (
	ClientUsername = "admin"
	ClientPassword = "W+JTd4bTfG.x8GxVORqQNK"
)

func Example_client() {
	packet := radius.New(radius.CodeAccessRequest, []byte(`PJbV6OoF9DvR1RHDTarTbil5pWIV/vbDnKbGgfx0krc=`))
	rfc2865.UserName_SetString(packet, ClientUsername)
	rfc2865.UserPassword_SetString(packet, ClientPassword)
	rfc2865.NASIdentifier_SetString(packet, "ngs0EJehk7x5MrFub8WbmaTIcP9atO5Ehc92mY8ptIB")
	client := &radius.Client{
		Retry:           time.Second,
		MaxPacketErrors: 1,
	}
	response, err := client.Exchange(context.Background(), packet, "localhost:18182")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Code:", response.Code)
}
func TestExample_client(t *testing.T) {
	Example_client()
}
