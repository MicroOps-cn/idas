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

package transport

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"

	"github.com/MicroOps-cn/fuck/log"
	"github.com/go-kit/log/level"
	uuid "github.com/satori/go.uuid"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2869"
)

type RadiusMessageAuthenticatorResponseWriter struct {
	w radius.ResponseWriter
}

func (r *RadiusMessageAuthenticatorResponseWriter) Write(packet *radius.Packet) error {
	if _, ok := packet.Lookup(rfc2869.MessageAuthenticator_Type); !ok {
		if err := rfc2869.MessageAuthenticator_Set(packet, make([]uint8, 16)); err != nil {
			return err
		}
		hash := hmac.New(md5.New, packet.Secret)
		encode, err := packet.Encode()
		if err != nil {
			return err
		}
		hash.Write(encode)
		packet.Set(rfc2869.MessageAuthenticator_Type, hash.Sum(nil))
	}
	return r.w.Write(packet)
}

func RadiusEAPAuthFilter(w radius.ResponseWriter, r *radius.Request, fc *radius.FilterChain) {
	logger := log.GetContextLogger(r.Context())
	if _, ok := r.Lookup(rfc2869.EAPMessage_Type); ok {
		msgAuth, ok := r.Lookup(rfc2869.MessageAuthenticator_Type)
		if !ok {
			level.Error(logger).Log("msg", "invalid request", "err", "EAPMessage exists but MessageAuthenticator attribute missing")
			w.Write(r.Response(radius.CodeAccessReject))
			return
		}

		if err := rfc2869.MessageAuthenticator_Set(r.Packet, make([]uint8, 16)); err != nil {
			level.Error(logger).Log("msg", "invalid request", "err", err)
			w.Write(r.Response(radius.CodeAccessReject))
			return
		}
		hash := hmac.New(md5.New, r.Packet.Secret)
		encode, err := r.Encode()
		if err != nil {
			level.Error(logger).Log("msg", "invalid request", "err", err)
			w.Write(r.Response(radius.CodeAccessReject))
			return
		}
		hash.Write(encode)
		if !bytes.Equal(hash.Sum(nil), msgAuth) {
			level.Error(logger).Log("msg", "invalid request", "err", "MessageAuthenticator is bad")
			w.Write(r.Response(radius.CodeAccessReject))
			return
		}
		r.Packet.Set(rfc2869.MessageAuthenticator_Type, msgAuth)
		w = &RadiusMessageAuthenticatorResponseWriter{w: w}
	}
	fc.ProcessFilter(w, r)
}

type EAPMessageCode uint8

const (
	EAPMessageCodeRequest    EAPMessageCode = 1
	EAPMessageCodeResponse   EAPMessageCode = 2
	EAPMessageCodeSuccess    EAPMessageCode = 3
	EAPMessageCodeFailure    EAPMessageCode = 4
	EAPMessageCodeComplete   EAPMessageCode = 5
	EAPMessageCodeOutOfOrder EAPMessageCode = 6
)

type EAPMessageID uint8

type EAPMessageType uint8

const (
	EAPMessageTypeIdentity        EAPMessageType = 1
	EAPMessageTypeNotification    EAPMessageType = 2
	EAPMessageTypeNak             EAPMessageType = 3
	EAPMessageTypeMD5Challenge    EAPMessageType = 4
	EAPMessageTypeOneTimePassword EAPMessageType = 5
	EAPMessageTypeGTC             EAPMessageType = 6
	EAPMessageTypeGeneric         EAPMessageType = 254
	EAPMessageTypeExpandable      EAPMessageType = 255
)

type EAPMessage struct {
	Code EAPMessageCode
	Id   byte
	Type EAPMessageType
	Date []byte
}

type RadiusEAPResponseWriter struct {
	w        radius.ResponseWriter
	identity byte
}

func (r *RadiusEAPResponseWriter) Write(packet *radius.Packet) error {
	if _, ok := packet.Lookup(rfc2869.EAPMessage_Type); !ok {
		packet.Set(rfc2869.EAPMessage_Type, []byte{byte(EAPMessageCodeSuccess), r.identity, 0x00, 0x04})
	}
	return r.w.Write(packet)
}

func RadiusEAPHandler(writer radius.ResponseWriter, r *radius.Request, ch HandlerChain) {
	r.Context()
	if eapMessage, ok := r.Lookup(rfc2869.EAPMessage_Type); ok {
		msg := EAPMessage{
			Code: EAPMessageCode(eapMessage[0]),
			Id:   eapMessage[1],
			Type: EAPMessageType(eapMessage[4]),
			Date: eapMessage[5:],
		}
		writer = &RadiusEAPResponseWriter{
			w:        writer,
			identity: msg.Id,
		}
		switch msg.Code {
		case EAPMessageCodeResponse:
			switch msg.Type {
			case EAPMessageTypeIdentity:
				q := &radius.Packet{
					Code:       radius.CodeAccessChallenge,
					Identifier: r.Identifier,
					Secret:     r.Secret,
				}
				q.Attributes.Set(rfc2865.State_Type, uuid.Must(uuid.NewV4()).Bytes()[:16])
				copy(q.Authenticator[:], r.Authenticator[:])
				// EAP Type: GTC, EAP Data: []byte("Password:")
				q.Set(rfc2869.EAPMessage_Type, []byte{byte(EAPMessageCodeRequest), msg.Id, 0x00, 0x0f, byte(EAPMessageTypeGTC), 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x3a, 0x20})
				writer.Write(q)
				return
			case EAPMessageTypeGTC:
				r.Packet.Set(rfc2865.UserPassword_Type, msg.Date)
			default:
				level.Debug(log.GetContextLogger(r.Context())).Log("msg", "unknown EAP Message Type", "type", msg.Type)
				writer.Write(r.Response(radius.CodeAccessReject))
				return
			}
		default:
			level.Debug(log.GetContextLogger(r.Context())).Log("msg", "unknown EAP Message Code", "code", msg.Code)
			writer.Write(r.Response(radius.CodeAccessReject))
			return
		}
	}
	ch.ServeRADIUS(writer, r)
}
