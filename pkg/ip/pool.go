package ip

import (
	"fmt"
	"github.com/sshota0809/loadbalancerip-mutator/pkg/logger"
	"net"
	"strings"
)

type IpAddr string

func (ip IpAddr) ToString() string {
	return string(ip)
}

type IsUsed bool

type IPPool struct {
	IPs map[IpAddr]IsUsed
}

func NewIpPool(pool string) (*IPPool, error) {
	// split string that is specified by option to cidrList
	cidrList := strings.Split(pool, ",")

	ips := map[IpAddr]IsUsed{}
	for _, cidr := range cidrList {
		ipList, err := generateIPList(cidr)
		if err != nil {
			return nil, err
		}

		logger.Log.Debug(fmt.Sprintf("IP list in pool: [%s] %s", cidr, strings.Join(ipList, ",")))
		// initialize map of IPs
		for _, ip := range ipList {
			ips[IpAddr(ip)] = false
		}
	}

	return &IPPool{IPs: ips}, nil
}

func generateIPList(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, &CidrFormatError{error: err}
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
