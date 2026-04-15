package util_test

import (
	"net/netip"
	"testing"

	"github.com/juanfont/headscale/pkg/util"
)

func TestNormalizePrefix(t *testing.T) {
	prefix := netip.MustParsePrefix("10.0.0.1/24")
	normalized := util.NormalizePrefix(prefix)
	expected := netip.MustParsePrefix("10.0.0.0/24")
	if normalized != expected {
		t.Errorf("expected %s, got %s", expected, normalized)
	}
}

func TestIPsFromPrefix(t *testing.T) {
	prefix := netip.MustParsePrefix("192.168.1.0/30")
	addrs, err := util.IPsFromPrefix(prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// /30 has 4 addresses: network, 2 hosts, broadcast
	// we skip the network address, so expect 3
	if len(addrs) != 3 {
		t.Errorf("expected 3 addresses, got %d", len(addrs))
	}
}

func TestIPsFromPrefixInvalid(t *testing.T) {
	var prefix netip.Prefix
	_, err := util.IPsFromPrefix(prefix)
	if err == nil {
		t.Error("expected error for invalid prefix, got nil")
	}
}

func TestIsIPv4(t *testing.T) {
	addr := netip.MustParseAddr("10.0.0.1")
	if !util.IsIPv4(addr) {
		t.Error("expected IPv4 address to return true")
	}
	if util.IsIPv6(addr) {
		t.Error("expected IPv4 address to not be IPv6")
	}
}

func TestIsIPv6(t *testing.T) {
	addr := netip.MustParseAddr("fd7a:115c::1")
	if !util.IsIPv6(addr) {
		t.Error("expected IPv6 address to return true")
	}
	if util.IsIPv4(addr) {
		t.Error("expected IPv6 address to not be IPv4")
	}
}

func TestContainsAddr(t *testing.T) {
	prefix := netip.MustParsePrefix("10.0.0.0/24")
	inside := netip.MustParseAddr("10.0.0.42")
	outside := netip.MustParseAddr("10.0.1.1")

	if !util.ContainsAddr(prefix, inside) {
		t.Errorf("expected prefix to contain %s", inside)
	}
	if util.ContainsAddr(prefix, outside) {
		t.Errorf("expected prefix to not contain %s", outside)
	}
}
