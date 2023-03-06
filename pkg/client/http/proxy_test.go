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

package http

import (
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/stretchr/testify/require"
)

func TestSendProxyRequest(t *testing.T) {
	logs.SetDefaultLogger(logs.New(logs.WithConfig(logs.MustNewConfig("debug", "logfmt"))))
	serv := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("ok"))
		require.NoError(t, err)
	}))
	type args struct {
		host               string
		insecureSkipVerify bool
		caCert             *x509.CertPool
	}
	cert := serv.Certificate()
	certPool := x509.NewCertPool()
	var host string
	if len(cert.DNSNames) > 0 {
		host = cert.DNSNames[0]
	} else if len(cert.IPAddresses) > 0 {
		host = cert.IPAddresses[0].String()
	}
	certPool.AddCert(cert)
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr string
	}{
		{name: "bad certificate", want: "", wantErr: "certificate signed by unknown authority"},
		{name: "unknown authority", args: args{host: host}, want: "", wantErr: "certificate signed by unknown authority"},
		{name: "ok", args: args{insecureSkipVerify: true}, want: "ok", wantErr: ""},
		{name: "ok", args: args{caCert: certPool, host: "xxxxxxx"}, want: "ok", wantErr: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", serv.URL, nil)
			require.NoError(t, err)

			if tt.args.insecureSkipVerify {
				req = WithInsecureSkipVerify(req)
			}
			if tt.args.caCert != nil {
				req = WithCaCert(req, tt.args.caCert)
			}
			if len(tt.args.host) > 0 {
				req.Host = tt.args.host
			}
			got, err := SendProxyRequest(req)
			if (err == nil && len(tt.wantErr) > 0) || (err != nil && !strings.Contains(err.Error(), tt.wantErr)) {
				t.Errorf("SendProxyRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}
			all, err := ioutil.ReadAll(got.Body)
			require.NoError(t, err)
			if !reflect.DeepEqual(string(all), tt.want) {
				t.Errorf("SendProxyRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
