package netipz

import (
	"errors"
	"fmt"
	"net/netip"
	"slices"

	"github.com/hakadoriya/z.go/mustz"
)

//nolint:gochecknoglobals
var (
	Zero                   = netip.Addr{}
	LoopbackAddress        = mustz.One(netip.ParsePrefix("127.0.0.0/8"))
	LinkLocalAddress       = mustz.One(netip.ParsePrefix("169.254.0.0/16"))
	PrivateIPAddressClassA = mustz.One(netip.ParsePrefix("10.0.0.0/8"))
	PrivateIPAddressClassB = mustz.One(netip.ParsePrefix("172.16.0.0/12"))
	PrivateIPAddressClassC = mustz.One(netip.ParsePrefix("192.168.0.0/16"))
)

var ErrPrefixesIsEmpty = errors.New("cidrs is empty")

func MustParsePrefixes(prefixes ...string) []netip.Prefix {
	ipNets, err := ParsePrefixes(prefixes...)
	if err != nil {
		panic(fmt.Errorf("ParsePrefixes: %w", err))
	}
	return ipNets
}

func ParsePrefixes(prefixes ...string) ([]netip.Prefix, error) {
	if len(prefixes) == 0 {
		return nil, fmt.Errorf("prefixes=%v: %w", prefixes, ErrPrefixesIsEmpty)
	}

	ipNets := make([]netip.Prefix, len(prefixes))
	for idx, prefix := range prefixes {
		ipNet, err := netip.ParsePrefix(prefix)
		if err != nil {
			return nil, fmt.Errorf("ParsePrefix: prefix=%v: %w", prefix, err)
		}
		ipNets[idx] = ipNet
	}
	return ipNets, nil
}

type IPNetSet []netip.Prefix

func (ipNetSet IPNetSet) Contains(ip netip.Addr) bool {
	return slices.ContainsFunc(ipNetSet, func(ipNet netip.Prefix) bool {
		return ipNet.Contains(ip)
	})
}
