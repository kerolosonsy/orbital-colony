package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"time"
)

type ColonyTLSOptions struct {
	Server      tls.Certificate
	Trust       *x509.CertPool
	ClientTrust *x509.CertPool
	Client      *tls.Certificate
	ServerName  string
	Now         func() time.Time
}

func ColonyTLS(opts ColonyTLSOptions) (*tls.Config, *tls.Config, error) {
	if len(opts.Server.Certificate) == 0 || opts.Server.PrivateKey == nil || opts.Trust == nil || opts.ServerName == "" || opts.Now == nil {
		return nil, nil, errors.New("invalid TLS options")
	}
	server := &tls.Config{Certificates: []tls.Certificate{opts.Server}, MinVersion: tls.VersionTLS13, Time: opts.Now}
	client := &tls.Config{RootCAs: opts.Trust.Clone(), ServerName: opts.ServerName, MinVersion: tls.VersionTLS13, Time: opts.Now}
	if opts.ClientTrust != nil {
		server.ClientCAs = opts.ClientTrust.Clone()
		server.ClientAuth = tls.RequireAndVerifyClientCert
	}
	if opts.Client != nil {
		client.Certificates = []tls.Certificate{*opts.Client}
	}
	return server, client, nil
}
func main() {}
