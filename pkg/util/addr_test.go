package util_test

import (
	"net/netip"
	"testing"

	"github.com/juanfont/headscale/pkg/util"
)

func TestParseAddr(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid IPv4", "192.168.1.1", false},
		{"valid IPv6", "fd7a:115c:a1e0::1", false},
		{"invalid", "not-an-ip", true},
		{"empty", "", true},
		{"CIDR notation", "10.0.0.1/24", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := util.ParseAddr(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAddr(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestParseAddrs(t *testing.T) {
	inputs := []string{"10.0.0.1", "10.0.0.2", "fd7a::1"}
	addrs, err := util.ParseAddrs(inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(addrs) != len(inputs) {
		t.Errorf("expected %d addrs, got %d", len(inputs), len(addrs))
	}

	_, err = util.ParseAddrs([]string{"bad"})
	if err == nil {
		t.Error("expected error for invalid addr")
	}
}

func TestAddrFamily(t *testing.T) {
	if fam := util.AddrFamily(netip.MustParseAddr("10.0.0.1")); fam != 4 {
		t.Errorf("expected family 4, got %d", fam)
	}
	if fam := util.AddrFamily(netip.MustParseAddr("fd7a::1")); fam != 6 {
		t.Errorf("expected family 6, got %d", fam)
	}
}

func TestAddrSetContains(t *testing.T) {
	set := []netip.Addr{
		netip.MustParseAddr("10.0.0.1"),
		netip.MustParseAddr("10.0.0.2"),
	}
	if !util.AddrSetContains(set, netip.MustParseAddr("10.0.0.1")) {
		t.Error("expected set to contain 10.0.0.1")
	}
	if util.AddrSetContains(set, netip.MustParseAddr("10.0.0.3")) {
		t.Error("expected set not to contain 10.0.0.3")
	}
}

func TestUniqueAddrs(t *testing.T) {
	addrs := []netip.Addr{
		netip.MustParseAddr("10.0.0.1"),
		netip.MustParseAddr("10.0.0.1"),
		netip.MustParseAddr("10.0.0.2"),
	}
	uniq := util.UniqueAddrs(addrs)
	if len(uniq) != 2 {
		t.Errorf("expected 2 unique addrs, got %d", len(uniq))
	}
}
