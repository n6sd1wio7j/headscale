package util

import (
	"net/netip"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOverlapsPrefix(t *testing.T) {
	tests := []struct {
		name  string
		a, b  string
		want  bool
	}{
		{"same prefix", "10.0.0.0/8", "10.0.0.0/8", true},
		{"overlapping", "10.0.0.0/8", "10.1.0.0/16", true},
		{"non-overlapping", "10.0.0.0/8", "192.168.0.0/16", false},
		{"ipv6 overlap", "fd7a::/16", "fd7a:115c::/32", true},
		{"ipv6 no overlap", "fd7a::/16", "fd00::/16", false},
		// Edge case: adjacent but non-overlapping subnets
		{"adjacent subnets", "10.0.0.0/25", "10.0.0.128/25", false},
		// Edge case: /0 (default route) overlaps everything
		{"default route overlaps all", "0.0.0.0/0", "192.168.1.0/24", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := netip.MustParsePrefix(tt.a)
			b := netip.MustParsePrefix(tt.b)
			got := OverlapsPrefix(a, b)
			if got != tt.want {
				t.Errorf("OverlapsPrefix(%s, %s) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestPrefixesOverlap(t *testing.T) {
	tests := []struct {
		name     string
		prefixes []string
		want     bool
	}{
		{"no overlap", []string{"10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12"}, false},
		{"overlap", []string{"10.0.0.0/8", "10.1.0.0/16"}, true},
		{"single", []string{"10.0.0.0/8"}, false},
		{"empty", []string{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfxs, _ := ParsePrefixes(tt.prefixes)
			got := PrefixesOverlap(pfxs)
			if got != tt.want {
				t.Errorf("PrefixesOverlap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePrefixes(t *testing.T) {
	got, err := ParsePrefixes([]string{"10.1.2.3/8", "192.168.1.100/24"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []netip.Prefix{
		netip.MustParsePrefix("10.0.0.0/8"),
		netip.MustParsePrefix("192.168.1.0/24"),
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ParsePrefixes() mismatch (-want +got):\n%s", diff)
	}

	_, err = ParsePrefixes([]string{"not-a-cidr"})
	if err == nil {
		t.Error("expected error for invalid CIDR, got nil")
	}
}

func TestContainsPrefix(t *testing.T) {
	tests := []struct {
		name        string
		outer, inner string
		want        bool
	}{
		{"contains", "10.0.0.0/8", "10.1.0.0/16", true},
		{"equal", "10.0.0.0/8", "10.0.0.0/8", true},
		{"not contains", "10.1.0.0/16", "10.0.0.0/8", false},
		{"disjoint", "10.0.0.0/8", "192.168.0.0/16", false},
		// Verify that a /32 host route is contained within its parent network
		{"host route contained", "10.0.0.0/24", "10.0.0.1/32", true},
		// Verify that a /128 IPv6 host route is contained within its parent network
		{"ipv6 host route contained", "fd7a:115c:a1e0::/48", "fd7a:115c:a1e0::1/128", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outer := netip.MustParsePrefix(tt.outer)
			inner := netip.MustParsePrefix(tt.inner)
			got := ContainsPrefix(outer, inner)
			if got != tt.want {
				t.Errorf("ContainsPrefix(%s, %s) = %v, want %v", tt.outer, tt.inner, got, tt.want)
			}
		})
	}
}
