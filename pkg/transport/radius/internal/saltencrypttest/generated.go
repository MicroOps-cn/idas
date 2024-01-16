// Code generated by radius-dict-gen. DO NOT EDIT.

package saltencrypttest

import (
	"crypto/rand"
	"net"
	"strconv"

	"layeh.com/radius"
)

const (
	SEInteger_Type radius.Type = 50
	SEOctets_Type  radius.Type = 51
	SEIPAddr_Type  radius.Type = 52
)

type SEInteger uint32

var SEInteger_Strings = map[SEInteger]string{}

func (a SEInteger) String() string {
	if str, ok := SEInteger_Strings[a]; ok {
		return str
	}
	return "SEInteger(" + strconv.FormatUint(uint64(a), 10) + ")"
}

func SEInteger_Add(p *radius.Packet, value SEInteger) (err error) {
	a := radius.NewInteger(uint32(value))
	var salt [2]byte
	_, err = rand.Read(salt[:])
	if err != nil {
		return
	}
	salt[0] |= 1 << 7
	a, err = radius.NewTunnelPassword(a, salt[:], p.Secret, p.Authenticator[:])
	p.Add(SEInteger_Type, a)
	return
}

func SEInteger_Get(p, q *radius.Packet) (value SEInteger) {
	value, _ = SEInteger_Lookup(p, q)
	return
}

func SEInteger_Gets(p, q *radius.Packet) (values []SEInteger, err error) {
	var i uint32
	for _, avp := range p.Attributes {
		if avp.Type != SEInteger_Type {
			continue
		}
		attr := avp.Attribute
		attr, _, err = radius.TunnelPassword(attr, p.Secret, q.Authenticator[:])
		if err != nil {
			return
		}
		i, err = radius.Integer(attr)
		if err != nil {
			return
		}
		values = append(values, SEInteger(i))
	}
	return
}

func SEInteger_Lookup(p, q *radius.Packet) (value SEInteger, err error) {
	a, ok := p.Lookup(SEInteger_Type)
	if !ok {
		err = radius.ErrNoAttribute
		return
	}
	a, _, err = radius.TunnelPassword(a, p.Secret, q.Authenticator[:])
	if err != nil {
		return
	}
	var i uint32
	i, err = radius.Integer(a)
	if err != nil {
		return
	}
	value = SEInteger(i)
	return
}

func SEInteger_Set(p *radius.Packet, value SEInteger) (err error) {
	a := radius.NewInteger(uint32(value))
	p.Set(SEInteger_Type, a)
	return
}

func SEInteger_Del(p *radius.Packet) {
	p.Attributes.Del(SEInteger_Type)
}

func SEOctets_Add(p *radius.Packet, value []byte) (err error) {
	var a radius.Attribute
	var salt [2]byte
	_, err = rand.Read(salt[:])
	if err != nil {
		return
	}
	salt[0] |= 1 << 7
	a, err = radius.NewTunnelPassword(value, salt[:], p.Secret, p.Authenticator[:])
	if err != nil {
		return
	}
	p.Add(SEOctets_Type, a)
	return
}

func SEOctets_AddString(p *radius.Packet, value string) (err error) {
	var a radius.Attribute
	var salt [2]byte
	_, err = rand.Read(salt[:])
	if err != nil {
		return
	}
	salt[0] |= 1 << 7
	a, err = radius.NewTunnelPassword([]byte(value), salt[:], p.Secret, p.Authenticator[:])
	if err != nil {
		return
	}
	p.Add(SEOctets_Type, a)
	return
}

func SEOctets_Get(p, q *radius.Packet) (value []byte) {
	value, _ = SEOctets_Lookup(p, q)
	return
}

func SEOctets_GetString(p, q *radius.Packet) (value string) {
	value, _ = SEOctets_LookupString(p, q)
	return
}

func SEOctets_Gets(p, q *radius.Packet) (values [][]byte, err error) {
	var i []byte
	for _, avp := range p.Attributes {
		if avp.Type != SEOctets_Type {
			continue
		}
		attr := avp.Attribute
		i, _, err = radius.TunnelPassword(attr, p.Secret, q.Authenticator[:])
		if err != nil {
			return
		}
		values = append(values, i)
	}
	return
}

func SEOctets_GetStrings(p, q *radius.Packet) (values []string, err error) {
	var i string
	for _, avp := range p.Attributes {
		if avp.Type != SEOctets_Type {
			continue
		}
		attr := avp.Attribute
		var up []byte
		up, _, err = radius.TunnelPassword(attr, p.Secret, q.Authenticator[:])
		if err == nil {
			i = string(up)
		}
		if err != nil {
			return
		}
		values = append(values, i)
	}
	return
}

func SEOctets_Lookup(p, q *radius.Packet) (value []byte, err error) {
	a, ok := p.Lookup(SEOctets_Type)
	if !ok {
		err = radius.ErrNoAttribute
		return
	}
	value, _, err = radius.TunnelPassword(a, p.Secret, q.Authenticator[:])
	return
}

func SEOctets_LookupString(p, q *radius.Packet) (value string, err error) {
	a, ok := p.Lookup(SEOctets_Type)
	if !ok {
		err = radius.ErrNoAttribute
		return
	}
	var b []byte
	b, _, err = radius.TunnelPassword(a, p.Secret, q.Authenticator[:])
	if err == nil {
		value = string(b)
	}
	return
}

func SEOctets_Set(p *radius.Packet, value []byte) (err error) {
	var a radius.Attribute
	var salt [2]byte
	_, err = rand.Read(salt[:])
	if err != nil {
		return
	}
	salt[0] |= 1 << 7
	a, err = radius.NewTunnelPassword(value, salt[:], p.Secret, p.Authenticator[:])
	if err != nil {
		return
	}
	p.Set(SEOctets_Type, a)
	return
}

func SEOctets_SetString(p *radius.Packet, value string) (err error) {
	var a radius.Attribute
	var salt [2]byte
	_, err = rand.Read(salt[:])
	if err != nil {
		return
	}
	salt[0] |= 1 << 7
	a, err = radius.NewTunnelPassword([]byte(value), salt[:], p.Secret, p.Authenticator[:])
	if err != nil {
		return
	}
	p.Set(SEOctets_Type, a)
	return
}

func SEOctets_Del(p *radius.Packet) {
	p.Attributes.Del(SEOctets_Type)
}

func SEIPAddr_Add(p *radius.Packet, value net.IP) (err error) {
	var a radius.Attribute
	a, err = radius.NewIPAddr(value)
	if err != nil {
		return
	}
	var salt [2]byte
	_, err = rand.Read(salt[:])
	if err != nil {
		return
	}
	salt[0] |= 1 << 7
	a, err = radius.NewTunnelPassword(a, salt[:], p.Secret, p.Authenticator[:])
	p.Add(SEIPAddr_Type, a)
	return
}

func SEIPAddr_Get(p, q *radius.Packet) (value net.IP) {
	value, _ = SEIPAddr_Lookup(p, q)
	return
}

func SEIPAddr_Gets(p, q *radius.Packet) (values []net.IP, err error) {
	var i net.IP
	for _, avp := range p.Attributes {
		if avp.Type != SEIPAddr_Type {
			continue
		}
		attr := avp.Attribute
		attr, _, err = radius.TunnelPassword(attr, p.Secret, q.Authenticator[:])
		if err != nil {
			return
		}
		i, err = radius.IPAddr(attr)
		if err != nil {
			return
		}
		values = append(values, i)
	}
	return
}

func SEIPAddr_Lookup(p, q *radius.Packet) (value net.IP, err error) {
	a, ok := p.Lookup(SEIPAddr_Type)
	if !ok {
		err = radius.ErrNoAttribute
		return
	}
	a, _, err = radius.TunnelPassword(a, p.Secret, q.Authenticator[:])
	if err != nil {
		return
	}
	value, err = radius.IPAddr(a)
	return
}

func SEIPAddr_Set(p *radius.Packet, value net.IP) (err error) {
	var a radius.Attribute
	a, err = radius.NewIPAddr(value)
	if err != nil {
		return
	}
	var salt [2]byte
	_, err = rand.Read(salt[:])
	if err != nil {
		return
	}
	salt[0] |= 1 << 7
	a, err = radius.NewTunnelPassword(a, salt[:], p.Secret, p.Authenticator[:])
	p.Set(SEIPAddr_Type, a)
	return
}

func SEIPAddr_Del(p *radius.Packet) {
	p.Attributes.Del(SEIPAddr_Type)
}
