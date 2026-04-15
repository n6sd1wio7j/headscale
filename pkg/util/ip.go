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
func IPsFromPrefix(prefix netip.Prefix) ([]netip.Addr, error) {
	if !prefix.IsValid() {
		return nil, fmt.Errorf("invalid prefix: %s", prefix)
	}

	var addrs []netip.Addr
	addr := prefix.Addr()
	for {
		if !prefix.Contains(addr) {
			break
		}
		if addr != prefix.Masked().Addr() {
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
