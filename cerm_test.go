package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

//
// func signPEM() {
// 	var x509ca *x509.Certificate
// 	tart := time.Unix(0, 0)
// 	end, err := time.Parse("2006-01-02", "2049-12-31")
// 	if err != nil {
// 		panic(err)
// 	}
// 	template := x509.Certificate{
// 		SerialNumber: serial,
// 		Issuer:       x509ca.Subject,
// 		Subject: pkix.Name{
// 			Organization: []string{"Able untrusted MITM proxy Inc"},
// 		},
// 		NotBefore: start,
// 		NotAfter:  end,
//
// 		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
// 		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
// 		BasicConstraintsValid: true,
// 	}
//
// }

func Test_loadca(t *testing.T) {
	// root 证书
	rootCA, err := tls.LoadX509KeyPair("./openssl/myCA.cer", "./openssl/myCA.key")
	if err != nil {
		t.Error(err)
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

	template.DNSNames = append(template.DNSNames, "baidu.com")
	template.DNSNames = append(template.DNSNames, "*.baidu.com")
	template.Subject.CommonName = "*.baidu.com"
	// server.der key
	rootCert, err := x509.ParseCertificate(rootCA.Certificate[0])
	if err != nil {
		panic(err)
	}
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, rootCert, &pk.PublicKey, rootCA.PrivateKey)

	//INFO:pem.Encode(tofile, xxxx)
	certBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyServer := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	cert, err := tls.X509KeyPair(certBytes, keyServer)
	if err != nil {
		t.Error(err)
	}
	t.Log(cert)
}
