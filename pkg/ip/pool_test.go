package ip

import (
	"github.com/sshota0809/loadbalancerip-mutator/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initLogger() {
	logger.Init("error")
}

func TestNewIpPool(t *testing.T) {
	tests := []struct {
		description string
		pool        string
		expIPList   []IpAddr
		expErr      error
	}{
		{
			description: "Valid format single CIDR",
			pool:        "10.10.10.240/29",
			expIPList: []IpAddr{
				"10.10.10.240",
				"10.10.10.241",
				"10.10.10.242",
				"10.10.10.243",
				"10.10.10.244",
				"10.10.10.245",
				"10.10.10.246",
				"10.10.10.247",
			},
			expErr: nil,
		},
		{
			description: "Invalid format single CIDR",
			pool:        "10.10.10.240/33",
			expErr:      &CidrFormatError{},
		},
		{
			description: "Valid format multi CIDR",
			pool:        "10.10.10.240/29,10.10.10.10/32",
			expIPList: []IpAddr{
				"10.10.10.10",
				"10.10.10.240",
				"10.10.10.241",
				"10.10.10.242",
				"10.10.10.243",
				"10.10.10.244",
				"10.10.10.245",
				"10.10.10.246",
				"10.10.10.247",
			},
			expErr: nil,
		},
		{
			description: "Invalid format multi CIDR",
			pool:        "10.10.10.240/32,10.10.10.241/34",
			expErr:      &CidrFormatError{},
		},
	}

	initLogger()

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			pool, err := NewIpPool(tt.pool)

			if err != nil {
				assert.ErrorAs(t, err, &tt.expErr, "Error should be equal to expErr")
				return
			}

			var ipList []IpAddr
			for ipAddr, _ := range pool.IPs {
				ipList = append(ipList, ipAddr)
			}
			assert.ElementsMatch(t, ipList, tt.expIPList, "IPs should be equal to expIPList")
		})
	}
}
