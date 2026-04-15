package util

import (
	"fmt"
	"net/netip"
)

// NormalizePrefix ensures the prefix is in canonical form.
func NormalizePrefix(prefix netip.Prefix) netip.Prefix {
	return prefix.Masked()
}

// IPsFromPrefix returns all usable host addresses in a given prefix.
// Note: for IPv4, this excludes the network address and broadcast address.
// For IPv6, there is no broadcast address, so only the network address is excluded.
func IPsFromPrefix(prefix netip.Prefix) ([]netip.Addr, error) {
	if !prefix.IsValid() {
		return nil, fmt.Errorf("invalid prefix: %s", prefix)
	}

	networkAddr := prefix.Masked().Addr()

	var addrs []netip.Addr
	addr := prefix.Addr()
	for {
		if !prefix.Contains(addr) {
			break
		}
		if addr != networkAddr {
			addrs = append(addrs, addr)
		}
		addr = addr.Next()
	}

	return addrs, nil
}

// IsIPv4 returns true if the given address is an IPv4 address.
func IsIPv4(addr netip.Addr) bool {
	return addr.Is4()
}

// IsIPv6 returns true if the given address is an IPv6 address.
func IsIPv6(addr netip.Addr) bool {
	return addr.Is6() && !addr.Is4In6()
}

// ContainsAddr checks whether the given prefix contains the given address.
func ContainsAddr(prefix netip.Prefix, addr netip.Addr) bool {
	return prefix.Contains(addr)
}
