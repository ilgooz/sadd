// Package sadd helps you to parse multiple service addresses formatted in a special single string syntax
// specified atParseQuery()
package sadd

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var invalidAddressError = errors.New("Addres must be formatted as [?ip/host]:[0<port<=65535]")
var invalidAddressRangeError = errors.New(`Address range must be formatted as
[?ip/host]:[0<port<=65535]-[?ip2/host2]:[0<port2<=65535] ip1/host1<=ip2/host2 and port<port2`)

// Address is a service address.
type Address struct {
	Host string
	Port int
}

// Parse parses an address as [?ip/host]:[0<port<=65535]
// examples:
//  - localhost
//  - 192.168.1.22
//  - :3000
//  - localhost:3000
//  - 192.168.1.22:3000
func Parse(s string) (address *Address, err error) {
	address = &Address{}

	if strings.Contains(s, ":") {
		u := strings.Split(s, ":")
		if len(u) < 1 || len(u) > 2 {
			return nil, addressError{s}
		}

		address.Host = u[0]
		port := u[1]

		if port == "" {
			return nil, addressError{s}
		} else {
			i, err := strconv.Atoi(port)
			if err != nil {
				return nil, addressError{port}
			}

			switch {
			case i < 0:
				return nil, addressError{port}
			case i > 65535:
				return nil, addressError{port}
			}
			address.Port = i
		}
	} else {
		address.Host = s
		address.Port = 80
	}

	return address, err
}

// String formats an Address as string.
func (address *Address) String() string {
	return fmt.Sprintf("%s:%d", address.Host, address.Port)
}

// ParseQuery parses address query. Addresses in queries treated same as Parse but
// also supports comma sparated addresses and port ranges.
//
// examples:
//  - :3000
//  - :3000-:3030
//  - 192.168.1.22:5000-:5050
//  - 192.168.1.22:5000-192.168.1.22:5050
//  - 192.168.1.22:5000-192.168.1.33:5050
//  - :6379,192.168.1.22:5000,:3000-:3030,192.168.1.22:5000-192.168.1.33:5050
func ParseQuery(query string) ([]*Address, error) {
	qp := &queryParser{q: query}
	return qp.parse()
}

type queryParser struct {
	q string
}

func (qp *queryParser) parse() (addresses []*Address, err error) {
	sp := strings.Split(qp.q, ",")

	for _, v := range sp {
		if strings.Contains(v, "-") {
			sp2 := strings.Split(v, "-")
			if len(sp2) != 2 {
				return nil, addressRangeError{}
			}

			address1, err := Parse(sp2[0])
			if err != nil {
				return nil, err
			}
			address2, err := Parse(sp2[1])
			if err != nil {
				return nil, err
			}

			ip1 := net.ParseIP(address1.Host)
			ip2 := net.ParseIP(address2.Host)

			if ip2 == nil && address2.Host != "" && address1.Host != address2.Host {
				return nil, addressRangeError{address2.Host}
			}
			if address1.Port >= address2.Port {
				return nil, addressRangeError{fmt.Sprintf("%d, %d", address1.Port, address2.Port)}
			}

			var ips []string
			if ip1 != nil && ip2 != nil {
				if bytes.Compare(ip2, ip1) < 0 {
					return nil, addressRangeError{fmt.Sprintf("%s, %s", address1.Host, address2.Host)}
				}

				i := ip1.To4()
				ib := ip2.To4()
				v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
				vb := uint(ib[0])<<24 + uint(ib[1])<<16 + uint(ib[2])<<8 + uint(ib[3])

				for i := v; i <= vb; i++ {
					v3 := byte(v & 0xFF)
					v2 := byte((v >> 8) & 0xFF)
					v1 := byte((v >> 16) & 0xFF)
					v0 := byte((v >> 24) & 0xFF)
					ips = append(ips, net.IPv4(v0, v1, v2, v3).String())
					v += 1
				}
			}

			if len(ips) > 0 {
				for _, v := range ips {
					for i := address1.Port; i <= address2.Port; i++ {
						address := *address1
						address.Host = v
						address.Port = i
						addresses = append(addresses, &address)
					}
				}
			} else {
				for i := address1.Port; i <= address2.Port; i++ {
					address := *address1
					address.Port = i
					addresses = append(addresses, &address)
				}
			}

		} else {
			address, err := Parse(v)
			if err != nil {
				return nil, err
			}
			addresses = append(addresses, address)
		}
	}
	return addresses, nil
}

type addressError struct {
	s string
}

func (e addressError) Error() string {
	return fmt.Sprintf("Error while parsing `%s`: %s", e.s, invalidAddressError)
}

type addressRangeError struct {
	s string
}

func (e addressRangeError) Error() string {
	return fmt.Sprintf("Error while parsing `%s`: %s", e.s, invalidAddressRangeError)
}
