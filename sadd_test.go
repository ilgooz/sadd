package sadd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	var address *Address
	var err error

	address, err = Parse(":3000")
	assert.Nil(t, err)
	assert.Equal(t, "", address.Host)
	assert.Equal(t, 3000, address.Port)

	address, err = Parse("localhost:3000")
	assert.Nil(t, err)
	assert.Equal(t, "localhost", address.Host)
	assert.Equal(t, 3000, address.Port)

	address, err = Parse("192.168.1.126:3000")
	assert.Nil(t, err)
	assert.Equal(t, "192.168.1.126", address.Host)
	assert.Equal(t, 3000, address.Port)

	address, err = Parse("192.168.1.126")
	assert.Nil(t, err)
	assert.Equal(t, "192.168.1.126", address.Host)
	assert.Equal(t, 80, address.Port)

	address, err = Parse("localhost")
	assert.Nil(t, err)
	assert.Equal(t, "localhost", address.Host)
	assert.Equal(t, 80, address.Port)
}

func TestParseQuery(t *testing.T) {
	var addresses []*Address
	var err error

	addresses, err = ParseQuery(":3000")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(addresses))
	assert.Equal(t, "", addresses[0].Host)
	assert.Equal(t, 3000, addresses[0].Port)

	addresses, err = ParseQuery("localhost:3000")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(addresses))
	assert.Equal(t, "localhost", addresses[0].Host)
	assert.Equal(t, 3000, addresses[0].Port)

	addresses, err = ParseQuery("192.168.1.126:3000")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(addresses))
	assert.Equal(t, "192.168.1.126", addresses[0].Host)
	assert.Equal(t, 3000, addresses[0].Port)

	addresses, err = ParseQuery("192.168.1.126")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(addresses))
	assert.Equal(t, "192.168.1.126", addresses[0].Host)
	assert.Equal(t, 80, addresses[0].Port)

	addresses, err = ParseQuery("localhost")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(addresses))
	assert.Equal(t, "localhost", addresses[0].Host)
	assert.Equal(t, 80, addresses[0].Port)

	addresses, err = ParseQuery(":3000-:3003")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(addresses))
	assert.Equal(t, "", addresses[0].Host)
	assert.Equal(t, 3000, addresses[0].Port)
	assert.Equal(t, "", addresses[1].Host)
	assert.Equal(t, 3001, addresses[1].Port)
	assert.Equal(t, "", addresses[2].Host)
	assert.Equal(t, 3002, addresses[2].Port)
	assert.Equal(t, "", addresses[3].Host)
	assert.Equal(t, 3003, addresses[3].Port)

	addresses, err = ParseQuery(":6379,:3000-:3003,localhost:3000-:3003,192.168.1.126:3000-:3003")
	assert.Nil(t, err)
	assert.Equal(t, 13, len(addresses))
	assert.Equal(t, "", addresses[0].Host)
	assert.Equal(t, 6379, addresses[0].Port)
	assert.Equal(t, "", addresses[1].Host)
	assert.Equal(t, 3000, addresses[1].Port)
	assert.Equal(t, "", addresses[2].Host)
	assert.Equal(t, 3001, addresses[2].Port)
	assert.Equal(t, "", addresses[3].Host)
	assert.Equal(t, 3002, addresses[3].Port)
	assert.Equal(t, "", addresses[4].Host)
	assert.Equal(t, 3003, addresses[4].Port)
	assert.Equal(t, "localhost", addresses[5].Host)
	assert.Equal(t, 3000, addresses[5].Port)
	assert.Equal(t, "localhost", addresses[6].Host)
	assert.Equal(t, 3001, addresses[6].Port)
	assert.Equal(t, "localhost", addresses[7].Host)
	assert.Equal(t, 3002, addresses[7].Port)
	assert.Equal(t, "localhost", addresses[8].Host)
	assert.Equal(t, 3003, addresses[8].Port)
	assert.Equal(t, "192.168.1.126", addresses[9].Host)
	assert.Equal(t, 3000, addresses[9].Port)
	assert.Equal(t, "192.168.1.126", addresses[10].Host)
	assert.Equal(t, 3001, addresses[10].Port)
	assert.Equal(t, "192.168.1.126", addresses[11].Host)
	assert.Equal(t, 3002, addresses[11].Port)
	assert.Equal(t, "192.168.1.126", addresses[12].Host)
	assert.Equal(t, 3003, addresses[12].Port)

	_, err = ParseQuery("localhost-82")
	assert.Equal(t, addressRangeError{"82"}, err)

	_, err = ParseQuery("localhost:-82")
	assert.Equal(t, addressError{"localhost:"}, err)

	addresses, err = ParseQuery("localhost-:82")
	assert.Equal(t, 3, len(addresses))
	assert.Equal(t, "localhost", addresses[0].Host)
	assert.Equal(t, 80, addresses[0].Port)
	assert.Equal(t, "localhost", addresses[1].Host)
	assert.Equal(t, 81, addresses[1].Port)
	assert.Equal(t, "localhost", addresses[2].Host)
	assert.Equal(t, 82, addresses[2].Port)

	addresses, err = ParseQuery("localhost:80-:82")
	assert.Equal(t, 3, len(addresses))
	assert.Equal(t, "localhost", addresses[0].Host)
	assert.Equal(t, 80, addresses[0].Port)
	assert.Equal(t, "localhost", addresses[1].Host)
	assert.Equal(t, 81, addresses[1].Port)
	assert.Equal(t, "localhost", addresses[2].Host)
	assert.Equal(t, 82, addresses[2].Port)

	addresses, err = ParseQuery("localhost:80-localhost:81")
	assert.Equal(t, 2, len(addresses))
	assert.Equal(t, "localhost", addresses[0].Host)
	assert.Equal(t, 80, addresses[0].Port)
	assert.Equal(t, "localhost", addresses[1].Host)
	assert.Equal(t, 81, addresses[1].Port)

	_, err = ParseQuery(":3000-3003")
	assert.Equal(t, addressRangeError{"3003"}, err)

	_, err = ParseQuery(":3003-:3000")
	assert.Equal(t, addressRangeError{"3003, 3000"}, err)

	addresses, err = ParseQuery("192.168.1.26:3000-192.168.1.28:3001")
	assert.Equal(t, 6, len(addresses))
	assert.Equal(t, "192.168.1.26", addresses[0].Host)
	assert.Equal(t, 3000, addresses[0].Port)
	assert.Equal(t, "192.168.1.26", addresses[1].Host)
	assert.Equal(t, 3001, addresses[1].Port)
	assert.Equal(t, "192.168.1.27", addresses[2].Host)
	assert.Equal(t, 3000, addresses[2].Port)
	assert.Equal(t, "192.168.1.27", addresses[3].Host)
	assert.Equal(t, 3001, addresses[3].Port)
	assert.Equal(t, "192.168.1.28", addresses[4].Host)
	assert.Equal(t, 3000, addresses[4].Port)
	assert.Equal(t, "192.168.1.28", addresses[5].Host)
	assert.Equal(t, 3001, addresses[5].Port)

	_, err = ParseQuery("192.168.1.26:3000-192.168.1.25:3001")
	assert.Equal(t, addressRangeError{"192.168.1.26, 192.168.1.25"}, err)
}
