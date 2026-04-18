package util

import (
	"fmt"
	"net/netip"
)

// PrefixFamily returns the address family of the given prefix.
func PrefixFamily(prefix netip.Prefix) AddrFamily {
	return AddrFamily(prefix.Addr())
}

// PrefixHostAddr returns the host address from a prefix (the address
// without the mask applied), preserving the original address bits.
func PrefixHostAddr(prefix netip.Prefix) netip.Addr {
	return prefix.Addr()
}

// PrefixNetworkAddr returns the network address of the prefix
// (all host bits zeroed).
func PrefixNetworkAddr(prefix netip.Prefix) netip.Addr {
	return prefix.Masked().Addr()
}

// PrefixBroadcastAddr returns the broadcast address of the prefix
// (all host bits set). For IPv6 this is the last address in the range.
func PrefixBroadcastAddr(prefix netip.Prefix) netip.Addr {
	return lastAddr(prefix)
}

// PrefixSize returns the number of addresses in the prefix.
// For large prefixes (e.g. /0 IPv6) this may overflow uint64 and returns
// MaxUint64 as a sentinel value instead.
func PrefixSize(prefix netip.Prefix) uint64 {
	bits := prefix.Addr().BitLen() - prefix.Bits()
	if bits >= 64 {
		return ^uint64(0)
	}
	return 1 << uint(bits)
}

// SingleHostPrefix returns a /32 (IPv4) or /128 (IPv6) prefix for addr.
func SingleHostPrefix(addr netip.Addr) netip.Prefix {
	return netip.PrefixFrom(addr, addr.BitLen())
}

// PrefixesContainAddr reports whether any of the given prefixes contain addr.
func PrefixesContainAddr(prefixes []netip.Prefix, addr netip.Addr) bool {
	for _, p := range prefixes {
		if p.Contains(addr) {
			return true
		}
	}
	return false
}

// FilterPrefixesByFamily returns only the prefixes matching the given address family.
func FilterPrefixesByFamily(prefixes []netip.Prefix, family AddrFamily) []netip.Prefix {
	out := make([]netip.Prefix, 0, len(prefixes))
	for _, p := range prefixes {
		if AddrFamily(p.Addr()) == family {
			out = append(out, p)
		}
	}
	return out
}

// UniquePrefixes returns a deduplicated slice of prefixes preserving order.
func UniquePrefixes(prefixes []netip.Prefix) []netip.Prefix {
	seen := make(map[netip.Prefix]struct{}, len(prefixes))
	out := make([]netip.Prefix, 0, len(prefixes))
	for _, p := range prefixes {
		if _, ok := seen[p]; !ok {
			seen[p] = struct{}{}
			out = append(out, p)
		}
	}
	return out
}

// PrefixFromAddrAndBits constructs a prefix from an addr and bit length,
// returning an error if the combination is invalid.
// Note: addr.Prefix() already masks host bits, so the returned prefix may
// differ from the input addr if host bits are set.
func PrefixFromAddrAndBits(addr netip.Addr, bits int) (netip.Prefix, error) {
	p, err := addr.Prefix(bits)
	if err != nil {
		return netip.Prefix{}, fmt.Errorf("constructing prefix from %s/%d: %w", addr, bits, err)
	}
	return p, nil
}
