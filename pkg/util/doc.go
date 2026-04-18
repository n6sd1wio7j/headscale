// Package util provides shared utility functions used across the headscale
// codebase. It includes helpers for IP address and prefix manipulation,
// making it easier to work with Go's net/netip types throughout the project.
//
// # IP Utilities
//
// The ip.go file exposes functions for:
//   - Normalizing CIDR prefixes to their canonical (masked) form
//   - Enumerating host addresses within a prefix
//   - Distinguishing between IPv4 and IPv6 addresses
//   - Checking whether a prefix contains a given address
//
// Note: IPsFromPrefix can be expensive for large prefixes (e.g. /8 or larger).
// Avoid calling it in hot paths or on wide address ranges.
// As a rule of thumb, prefer prefixes of /24 or smaller when possible.
// For IPv6, be especially cautious — even a /64 contains 2^64 addresses.
// In practice, /112 or smaller is a safer upper bound for IPv6 enumeration.
//
// Example usage:
//
//	prefix := netip.MustParsePrefix("10.0.0.0/24")
//	addrs, err := util.IPsFromPrefix(prefix)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, addr := range addrs {
//		fmt.Println(addr)
//	}
package util
