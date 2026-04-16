// Package util provides shared utility functions for headscale.
package util

import (
	"fmt"
	"net/netip"
	"strings"
)

// ParseAddr parses a string into a netip.Addr, returning an error if the
// string is not a valid IP address.
func ParseAddr(s string) (netip.Addr, error) {
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("parsing IP address %q: %w", s, err)
	}
	return addr.Unmap(), nil
}

// ParseAddrs parses a slice of strings into a slice of netip.Addr values.
// Returns an error if any string is not a valid IP address.
func ParseAddrs(ss []string) ([]netip.Addr, error) {
	addrs := make([]netip.Addr, 0, len(ss))
	for _, s := range ss {
		addr, err := ParseAddr(s)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

// AddrFamily returns "IPv4" or "IPv6" depending on the address family of addr.
func AddrFamily(addr netip.Addr) string {
	if addr.Is4() {
		return "IPv4"
	}
	return "IPv6"
}

// AddrSetContains reports whether addrs contains addr.
func AddrSetContains(addrs []netip.Addr, addr netip.Addr) bool {
	for _, a := range addrs {
		if a == addr {
			return true
		}
	}
	return false
}

// UniqueAddrs returns a deduplicated copy of addrs, preserving order.
func UniqueAddrs(addrs []netip.Addr) []netip.Addr {
	seen := make(map[netip.Addr]struct{}, len(addrs))
	out := make([]netip.Addr, 0, len(addrs))
	for _, a := range addrs {
		if _, ok := seen[a]; !ok {
			seen[a] = struct{}{}
			out = append(out, a)
		}
	}
	return out
}

// AddrsToStrings converts a slice of netip.Addr to a slice of strings.
func AddrsToStrings(addrs []netip.Addr) []string {
	ss := make([]string, len(addrs))
	for i, a := range addrs {
		ss[i] = a.String()
	}
	return ss
}

// ParseAddrPort parses a host:port string into a netip.AddrPort.
// The host portion must be a valid IP address (not a hostname).
// Note: raw IPv6 addresses (e.g. "::1:8080") must be bracketed like "[::1]:8080".
func ParseAddrPort(s string) (netip.AddrPort, error) {
	// Handle bare IPv6 addresses without brackets by checking if it looks
	// like a raw IPv6 addr with no port.
	if strings.Count(s, ":") > 1 && !strings.HasPrefix(s, "[") {
		return netip.AddrPort{}, fmt.Errorf("parsing AddrPort %q: IPv6 address must be bracketed", s)
	}
	ap, err := netip.ParseAddrPort(s)
	if err != nil {
		return netip.AddrPort{}, fmt.Errorf("parsing AddrPort %q: %w", s, err)
	}
	return ap, nil
}
