package labkit

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"
)

// CA is a throwaway certificate authority for TLS lab tests.
type CA struct {
	Certificate *x509.Certificate
	key         *ecdsa.PrivateKey
}

// CertOptions describes a certificate to issue: names, validity and whether
// it is a client certificate.
type CertOptions struct {
	DNSNames    []string
	IPAddresses []net.IP
	NotBefore   time.Time
	NotAfter    time.Time
	Client      bool
}

// TestCA creates a new in-memory CA valid for ten years.
func TestCA() (*CA, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	template := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "Orbital lab CA"}, NotBefore: now.Add(-time.Hour), NotAfter: now.AddDate(10, 0, 0), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageCRLSign}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err
	}
	return &CA{Certificate: cert, key: key}, nil
}

// IssueCert can issue deliberately expired or mismatched certificates for tests.
func IssueCert(ca *CA, opts CertOptions) (tls.Certificate, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, err
	}
	if opts.NotBefore.IsZero() {
		opts.NotBefore = time.Now().Add(-time.Hour)
	}
	if opts.NotAfter.IsZero() {
		opts.NotAfter = time.Now().Add(time.Hour)
	}
	usage := x509.ExtKeyUsageServerAuth
	if opts.Client {
		usage = x509.ExtKeyUsageClientAuth
	}
	template := &x509.Certificate{SerialNumber: serial, DNSNames: opts.DNSNames, IPAddresses: opts.IPAddresses, NotBefore: opts.NotBefore, NotAfter: opts.NotAfter, KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{usage}}
	der, err := x509.CreateCertificate(rand.Reader, template, ca.Certificate, &key.PublicKey, ca.key)
	if err != nil {
		return tls.Certificate{}, err
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.Certificate{Certificate: [][]byte{der, ca.Certificate.Raw}, PrivateKey: key, Leaf: leaf}, nil
}
