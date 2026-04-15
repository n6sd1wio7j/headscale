package util

import (
	"fmt"
	"net/netip"
)

// OverlapsPrefix returns true if two prefixes overlap.
func OverlapsPrefix(a, b netip.Prefix) bool {
	return a.Overlaps(b)
}

// PrefixesOverlap checks if any prefix in the slice overlaps with any other.
func PrefixesOverlap(prefixes []netip.Prefix) bool {
	for i := 0; i < len(prefixes); i++ {
		for j := i + 1; j < len(prefixes); j++ {
			if prefixes[i].Overlaps(prefixes[j]) {
				return true
			}
		}
	}
	return false
}

// ParsePrefixes parses a slice of CIDR strings into netip.Prefix values.
func ParsePrefixes(cidrs []string) ([]netip.Prefix, error) {
	prefixes := make([]netip.Prefix, 0, len(cidrs))
	for _, cidr := range cidrs {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			return nil, fmt.Errorf("parsing prefix %q: %w", cidr, err)
		}
		prefixes = append(prefixes, NormalizePrefix(prefix))
	}
	return prefixes, nil
}

// ContainsPrefix returns true if outer fully contains inner.
func ContainsPrefix(outer, inner netip.Prefix) bool {
	if outer.Bits() > inner.Bits() {
		return false
	}
	return outer.Contains(inner.Addr()) && outer.Contains(lastAddr(inner))
}

// lastAddr returns the last address in a prefix.
func lastAddr(p netip.Prefix) netip.Addr {
	addr := p.Addr()
	bits := addr.BitLen()
	maskBits := p.Bits()

	a128 := addr.As16()
	for i := maskBits; i < bits; i++ {
		byteIdx := i / 8
		bitIdx := 7 - (i % 8)
		a128[byteIdx] |= 1 << uint(bitIdx)
	}
	if addr.Is4() {
		return netip.AddrFrom16(a128).Unmap()
	}
	return netip.AddrFrom16(a128)
}
