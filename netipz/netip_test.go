package netipz_test

import (
	"net/netip"
	"testing"

	"github.com/hakadoriya/z.go/errorz"
	"github.com/hakadoriya/z.go/netipz"
)

func TestParsePrefixes(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cidrs := []string{"10.0.0.0/8"}
		ipNets, err := netipz.ParsePrefixes(cidrs...)
		if err != nil {
			t.Errorf("❌: ParsePrefixes(%s) returned an error: %s", cidrs, err)
		}
		if actual := ipNets[0].String(); actual != cidrs[0] {
			t.Errorf("❌: ipNets[0].String(%s) != cidrs[0]: %s != %s", cidrs[0], actual, cidrs[0])
		}
	})
	t.Run("failure(FAILURE)", func(t *testing.T) {
		t.Parallel()

		cidrs := []string{"FAILURE"}
		if _, err := netipz.ParsePrefixes(cidrs...); err == nil {
			t.Errorf("❌: ParsePrefixes(%s) returned an error: %s", cidrs, err)
		}
	})
	t.Run("failure(empty)", func(t *testing.T) {
		t.Parallel()

		if _, err := netipz.ParsePrefixes(); err == nil {
			t.Errorf("❌: ParsePrefixes() returned an error: %s", err)
		}
	})
}

func TestMustParsePrefixes(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cidrs := []string{
			netipz.LoopbackAddress.String(),
			netipz.LinkLocalAddress.String(),
			netipz.PrivateIPAddressClassA.String(),
			netipz.PrivateIPAddressClassB.String(),
			netipz.PrivateIPAddressClassC.String(),
		}
		ipNets := netipz.MustParsePrefixes(cidrs...)
		for i := range ipNets {
			if expect, actual := cidrs[i], ipNets[i].String(); expect != actual {
				t.Fatalf("❌: MustParsePrefixes: expect(%s) != actual(%s)", cidrs, actual)
			}
		}
	})
	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		const cidr = "FAILURE"
		cidrs := []string{cidr}
		defer func() {
			err, ok := recover().(error)
			if !ok {
				t.Fatalf("❌: MustParsePrefixes should panic with an error")
			}
			if expect := `no '/'`; !errorz.Contains(err, expect) {
				t.Fatalf("❌: recover: expect(%s) != actual(%v)", expect, err)
			}
		}()
		netipz.MustParsePrefixes(cidrs...)
	})
}

func TestIPNetSet_Contains(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ipNetSet := netipz.IPNetSet([]netip.Prefix{netipz.PrivateIPAddressClassA})
		ip, err := netip.ParsePrefix("10.10.10.10/32")
		if err != nil {
			t.Fatalf("❌: netip.ParsePrefix: %s", err)
		}

		s := netipz.IPNetSet([]netip.Prefix{netipz.PrivateIPAddressClassA})

		if !s.Contains(ip.Addr()) {
			t.Errorf("❌: %s should contain %s", ipNetSet, ip)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		ipNetSet := netipz.IPNetSet([]netip.Prefix{netipz.PrivateIPAddressClassA})
		ip, err := netip.ParsePrefix("192.168.1.1/32")
		if err != nil {
			t.Fatalf("❌: netip.ParsePrefix: %s", err)
		}

		s := netipz.IPNetSet([]netip.Prefix{netipz.PrivateIPAddressClassA})

		if s.Contains(ip.Addr()) {
			t.Errorf("❌: %s should contain %s", ipNetSet, ip)
		}
	})
}
