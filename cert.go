package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"sync"
	"time"
)

var cert_host = map[string]tls.Certificate{}
var once = new(sync.Mutex)

func fetchCert(host string, ip string) (*tls.Certificate, error) {
	once.Lock()
	defer once.Unlock()
	if host == "" {
		host = ip
	}
	if cert, ok := cert_host[host]; ok {
		return &cert, nil
	} else {
		cert, err := genCert(host, ip)
		if err != nil {
			return nil, err
		}
		cert_host[host] = cert
		return &cert, err
	}
}

func genCert(host, ip string) (tls.Certificate, error) {
	// root 证书
	rootCA, err := tls.LoadX509KeyPair(config.RootCA, config.RootKey)
	if err != nil {
		panic(err)
	}
	// server 端key 生成
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	//  生成 服务器证书
	start := time.Unix(0, 0)
	end, _ := time.Parse("2006-01-02", "2049-12-31")
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Able untrusted MITM proxy Inc"},
		},
		NotBefore:             start,
		NotAfter:              end,
		IsCA:                  false,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	ipN := net.ParseIP(ip)

	template.IPAddresses = append(template.IPAddresses, ipN)
	template.DNSNames = append(template.DNSNames, host)
	template.Subject.CommonName = host
	// server.der key
	rootCert, err := x509.ParseCertificate(rootCA.Certificate[0])
	if err != nil {
		panic(err)
	}
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, rootCert, &pk.PublicKey, rootCA.PrivateKey)

	//INFO:pem.Encode(tofile, xxxx)
	certBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyServer := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	return tls.X509KeyPair(certBytes, keyServer)
}
