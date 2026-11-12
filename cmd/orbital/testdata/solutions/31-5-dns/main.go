package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"slices"
	"strings"
)

var ErrDockData = errors.New("invalid dock directory data")
var ErrDockAbsent = errors.New("dock name absent")

type DockResolver interface {
	Addresses(context.Context, string) ([]net.IP, error)
	Alias(context.Context, string) (string, error)
}
type AddressCard struct {
	Station, Canonical string
	A, AAAA            []string
}

func dockLabel(s string, max int) bool {
	if len(s) < 1 || len(s) > max || s[0] == '-' || s[len(s)-1] == '-' {
		return false
	}
	for _, b := range []byte(s) {
		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-') {
			return false
		}
	}
	return true
}
func dockCanonical(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	for _, b := range []byte(s) {
		if b >= 128 {
			return "", ErrDockData
		}
	}
	s = strings.TrimSuffix(strings.ToLower(s), ".")
	if len(s) > 253 {
		return "", ErrDockData
	}
	for _, label := range strings.Split(s, ".") {
		if !dockLabel(label, 63) {
			return "", ErrDockData
		}
	}
	return s, nil
}
func dockFailure(err error) error {
	var dns *net.DNSError
	if errors.As(err, &dns) && dns.IsNotFound {
		return fmt.Errorf("%w: %w", ErrDockAbsent, err)
	}
	return err
}
func dockAddresses(ips []net.IP) ([]string, []string, error) {
	a, aaaa := []string{}, []string{}
	seen := map[string]bool{}
	for _, ip := range ips {
		if ip.To16() == nil {
			return nil, nil, ErrDockData
		}
		text := ip.String()
		if seen[text] {
			continue
		}
		seen[text] = true
		if ip.To4() != nil {
			a = append(a, text)
		} else {
			aaaa = append(aaaa, text)
		}
	}
	slices.Sort(a)
	slices.Sort(aaaa)
	return a, aaaa, nil
}
func ResolveDock(ctx context.Context, resolver DockResolver, station string) (AddressCard, error) {
	if !dockLabel(station, 16) {
		return AddressCard{}, ErrDockData
	}
	if err := ctx.Err(); err != nil {
		return AddressCard{}, err
	}
	host := station + ".orbital.invalid"
	ips, err := resolver.Addresses(ctx, host)
	if ctx.Err() != nil {
		return AddressCard{}, ctx.Err()
	}
	if err != nil {
		return AddressCard{}, dockFailure(err)
	}
	alias, err := resolver.Alias(ctx, host)
	if ctx.Err() != nil {
		return AddressCard{}, ctx.Err()
	}
	if err != nil {
		return AddressCard{}, dockFailure(err)
	}
	canonical, err := dockCanonical(alias)
	if err != nil {
		return AddressCard{}, err
	}
	a, aaaa, err := dockAddresses(ips)
	if err != nil {
		return AddressCard{}, err
	}
	return AddressCard{Station: station, Canonical: canonical, A: a, AAAA: aaaa}, nil
}
func main() {}
