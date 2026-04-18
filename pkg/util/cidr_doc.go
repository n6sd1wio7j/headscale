// Package util provides network utility helpers for headscale.
//
// The cidr.go file extends the IP utilities with CIDR/prefix operations:
//
//   - OverlapsPrefix: checks whether two netip.Prefix values overlap.
//   - PrefixesOverlap: checks a slice of prefixes for any pairwise overlap;
//     useful for validating that configured IP pools do not conflict.
//   - ParsePrefixes: parses a slice of CIDR strings, normalising each one
//     (masking host bits) via NormalizePrefix before returning.
//   - ContainsPrefix: reports whether an outer prefix fully contains an
//     inner prefix, i.e. every address in inner is also in outer.
//
// These helpers are intended for use during configuration validation and
// address-allocation logic where overlapping or misconfigured ranges must
// be detected early.
//
// Note: ParsePrefixes silently drops invalid CIDR strings; callers should
// check that the returned slice length matches the input if strict
// validation is required.
package util
