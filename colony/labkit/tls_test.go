package labkit

import (
	"crypto/x509"
	"net"
	"testing"
	"time"
)

func TestIssuedCertificatesVerifyAndRejectInvalidIdentity(t *testing.T) {
	ca, err := TestCA()
	if err != nil {
		t.Fatal(err)
	}
	roots := x509.NewCertPool()
	roots.AddCert(ca.Certificate)
	for _, tc := range []struct {
		name  string
		opts  CertOptions
		dns   string
		usage x509.ExtKeyUsage
		ok    bool
	}{
		{"server", CertOptions{DNSNames: []string{"localhost"}}, "localhost", x509.ExtKeyUsageServerAuth, true},
		{"ip", CertOptions{IPAddresses: []net.IP{net.ParseIP("127.0.0.1")}}, "127.0.0.1", x509.ExtKeyUsageServerAuth, true},
		{"wrong-host", CertOptions{DNSNames: []string{"localhost"}}, "other", x509.ExtKeyUsageServerAuth, false},
		{"expired", CertOptions{DNSNames: []string{"localhost"}, NotBefore: time.Now().Add(-2 * time.Hour), NotAfter: time.Now().Add(-time.Hour)}, "localhost", x509.ExtKeyUsageServerAuth, false},
		{"client", CertOptions{Client: true}, "", x509.ExtKeyUsageClientAuth, true},
		{"wrong-usage", CertOptions{Client: true}, "", x509.ExtKeyUsageServerAuth, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cert, err := IssueCert(ca, tc.opts)
			if err != nil {
				t.Fatal(err)
			}
			_, err = cert.Leaf.Verify(x509.VerifyOptions{Roots: roots, DNSName: tc.dns, KeyUsages: []x509.ExtKeyUsage{tc.usage}})
			if (err == nil) != tc.ok {
				t.Fatal(err)
			}
			if cert.PrivateKey == nil || len(cert.Certificate) != 2 {
				t.Fatal("incomplete TLS certificate")
			}
		})
	}
}
